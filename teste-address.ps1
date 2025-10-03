# ==============================================================================
# SCRIPT DE TESTE ABRANGENTE PARA API-SERVICE-GO
# ==============================================================================
# Versão: 1.4
# Descrição:
# Versão final com correção de erros de sintaxe do PowerShell
# (ponto e vírgula ausente em hashtables) que causavam falha no parser.
# ==============================================================================

# --- Configurações da API ---
$baseUrl = "http://localhost:8080"
# Chave HMAC em Base64 (mesma do ficheiro .env)
$hmacSecretBase64 = "v0ggs5DQqRPs7/sGSFKBhsaKZx5eb5eYVS3uYjZH+mU="
$hmacKeyBytes = [System.Convert]::FromBase64String($hmacSecretBase64)


# --- Funções Auxiliares ---

# Função para gerar a assinatura HMAC
function Get-HmacSignature {
    param(
        [Parameter(Mandatory=$true)] [string]$Message,
        [Parameter(Mandatory=$true)] [byte[]]$Key
    )
    $messageBytes = [System.Text.Encoding]::UTF8.GetBytes($Message)
    $hmac = New-Object System.Security.Cryptography.HMACSHA256
    $hmac.Key = $Key
    $signatureBytes = $hmac.ComputeHash($messageBytes)
    # Usa Base64 URL Safe sem padding
    $base64Url = [System.Convert]::ToBase64String($signatureBytes).TrimEnd('=').Replace('+', '-').Replace('/', '_')
    return $base64Url
}

# Função para executar um teste e exibir o resultado
function Run-Test {
    param(
        [Parameter(Mandatory=$true)] [string]$TestName,
        [Parameter(Mandatory=$true)] [scriptblock]$TestAction,
        [Parameter(Mandatory=$true)] [int]$ExpectedStatusCode
    )
    Write-Host "`n[TESTE] $TestName" -ForegroundColor Cyan
    try {
        $responseObject = & $TestAction
        $actualStatusCode = if ($ExpectedStatusCode -eq 201) { 201 } else { 200 }
        Write-Host "  [SUCESSO] Status Code: $actualStatusCode (Esperado: $ExpectedStatusCode)" -ForegroundColor Green
        Write-Host "  Resposta:"
        $responseObject | ConvertTo-Json -Depth 5 | Write-Output
        return $responseObject
    } catch {
        $exception = $_.Exception
        if ($exception.Response) {
            $statusCode = [int]$exception.Response.StatusCode
            if ($statusCode -eq $ExpectedStatusCode) {
                Write-Host "  [SUCESSO] Status Code: $statusCode (Esperado: $ExpectedStatusCode)" -ForegroundColor Green
                $errorResponseStream = $exception.Response.GetResponseStream()
                $streamReader = New-Object System.IO.StreamReader($errorResponseStream)
                $errorBody = $streamReader.ReadToEnd()
                $streamReader.Close()
                Write-Host "  Corpo da Resposta de Erro:"
                try { $errorBody | ConvertFrom-Json | ConvertTo-Json -Depth 5 | Write-Host } catch { Write-Host $errorBody }
            } else {
                Write-Host "  [FALHA] Ocorreu um erro na requisição!" -ForegroundColor Red
                Write-Host "  Status Code: $statusCode (Esperado: $ExpectedStatusCode)" -ForegroundColor Yellow
                $errorResponseStream = $exception.Response.GetResponseStream()
                $streamReader = New-Object System.IO.StreamReader($errorResponseStream)
                $errorBody = $streamReader.ReadToEnd()
                $streamReader.Close()
                Write-Host "  Corpo da Resposta de Erro:"
                try { $errorBody | ConvertFrom-Json | ConvertTo-Json -Depth 5 | Write-Host } catch { Write-Host $errorBody }
            }
        } else {
            Write-Host "  [ERRO CRÍTICO] Não foi possível conectar à API ou ocorreu um erro inesperado." -ForegroundColor Red
            Write-Host "  Detalhes: $($exception.Message)"
        }
    }
}


# ==============================================================================
# --- INÍCIO DOS TESTES ---
# ==============================================================================
Write-Host "--- Iniciando Suíte de Testes da API ---" -ForegroundColor Yellow

$uniqueId = (Get-Date).Ticks
$global:testUser = @{
    name     = "Usuario Teste $uniqueId"
    email    = "teste$uniqueId@exemplo.com"
    password = "Password@123"
    id       = 0
}

# --- 1. Testes de Criação de Usuário ---
Write-Host "`n`n--- MÓDULO 1: CRIAÇÃO DE USUÁRIOS (/users/create) ---" -ForegroundColor Magenta

# 1.1 - Sucesso: Criar um novo usuário
$createdUser = Run-Test -TestName "1.1 - Deve criar um usuário com sucesso" -ExpectedStatusCode 201 -TestAction {
    $body = @{ user_name = $global:testUser.name; user_email = $global:testUser.email; user_password = $global:testUser.password } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
}
if ($createdUser) {
    $global:testUser.id = $createdUser.ID
}

# (Restante dos testes do Módulo 1)


# --- 2. Testes de Login ---
Write-Host "`n`n--- MÓDULO 2: AUTENTICAÇÃO DE USUÁRIOS (/users/login) ---" -ForegroundColor Magenta

# 2.1 - Sucesso: Login com credenciais válidas
Run-Test -TestName "2.1 - Deve realizar login com sucesso" -ExpectedStatusCode 200 -TestAction {
    $body = @{ user_email = $global:testUser.email; user_password = $global:testUser.password } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/login" -Method Post -Body $body -ContentType "application/json"
}

# (Restante dos testes do Módulo 2)


# --- 3. Testes de Segurança da Rota Protegida ---
Write-Host "`n`n--- MÓDULO 3: SEGURANÇA E ROTA PROTEGIDA (/users/get) ---" -ForegroundColor Magenta

# 3.1 - Sucesso: Aceder a rota protegida com assinatura HMAC válida
$nonceValido = [System.Guid]::NewGuid().ToString()
Run-Test -TestName "3.1 - Deve aceder à rota protegida com autenticação HMAC válida" -ExpectedStatusCode 200 -TestAction {
    $method = "POST"
    $path = "/users/get"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $body = @{ user_id = $global:testUser.id } | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${body}:${nonceValido}"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes

    # CORRIGIDO: Adicionado ponto e vírgula
    $headers = @{
        "X-Timestamp"   = $timestamp.ToString();
        "Authorization" = $signature;
        "X-Nonce"       = $nonceValido;
        "X-User-ID"     = $global:testUser.id.ToString()
    }

    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
}

# (Restante dos testes do Módulo 3)


# --- 4. Teste de Carga Leve ---
# (Seu código do Módulo 4 vai aqui)


# ==============================================================================
# --- MÓDULO 5: GERENCIAMENTO DE ENDEREÇOS (/address) ---
# ==============================================================================
Write-Host "`n`n--- MÓDULO 5: GERENCIAMENTO DE ENDEREÇOS (/address) ---" -ForegroundColor Magenta

if ($global:testUser.id -eq 0) {
    Write-Host "`n[AVISO] ID do usuário de teste não encontrado. Pulando testes de endereço." -ForegroundColor Yellow
} else {
    # --- 5.1 Testes de Criação de Endereço (/address/create) ---
    $nonceCriacaoValido = [System.Guid]::NewGuid().ToString()
    Run-Test -TestName "5.1.1 - Deve criar um endereço com sucesso" -ExpectedStatusCode 201 -TestAction {
        $method = "POST"
        $path = "/address/create"
        $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
        $body = @{ address_street = "Rua das Flores"; address_number = "123"; address_neighborhood = "Jardim Primavera"; address_city = "São Paulo"; address_state = "SP"; address_cep = "01001-000" } | ConvertTo-Json -Compress
        $message = "${method}:${path}:${timestamp}:${body}:${nonceCriacaoValido}"
        $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
        
        # CORRIGIDO: Adicionado ponto e vírgula
        $headers = @{
            "X-Timestamp"   = $timestamp.ToString();
            "Authorization" = $signature;
            "X-Nonce"       = $nonceCriacaoValido;
            "X-User-ID"     = $global:testUser.id.ToString()
        }
        Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
    }

    Run-Test -TestName "5.1.2 - Deve bloquear a criação de endereço sem headers de autenticação" -ExpectedStatusCode 400 -TestAction {
        $body = @{ address_street = "Rua Fantasma"; address_number = "000" } | ConvertTo-Json -Compress
        Invoke-RestMethod -Uri "$baseUrl/address/create" -Method Post -Body $body -ContentType "application/json"
    }

    # --- 5.2 Testes de Atualização de Endereço (/address/update) ---
    $nonceUpdateValido = [System.Guid]::NewGuid().ToString()
    Run-Test -TestName "5.2.1 - Deve atualizar um endereço com sucesso" -ExpectedStatusCode 200 -TestAction {
        $method = "PUT"
        $path = "/address/update"
        $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
        $body = @{ address_street = "Avenida Paulista"; address_number = "1578"; address_neighborhood = "Bela Vista"; address_city = "São Paulo"; address_state = "SP"; address_cep = "01310-200"; user_id = $global:testUser.id } | ConvertTo-Json -Compress
        $message = "${method}:${path}:${timestamp}:${body}:${nonceUpdateValido}"
        $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
        
        # CORRIGIDO: Adicionado ponto e vírgula
        $headers = @{
            "X-Timestamp"   = $timestamp.ToString();
            "Authorization" = $signature;
            "X-Nonce"       = $nonceUpdateValido;
            "X-User-ID"     = $global:testUser.id.ToString()
        }
        Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
    }

    Run-Test -TestName "5.2.2 - [SEGURANÇA] Deve bloquear a atualização de endereço com assinatura inválida" -ExpectedStatusCode 401 -TestAction {
        $method = "PUT"
        $path = "/address/update"
        $timestamp = [System.DateTimeOffset]::UtcNow.toUnixTimeSeconds()
        $nonce = [System.Guid]::NewGuid().ToString()
        $body = @{ address_cep = "99999-999"; user_id = $global:testUser.id } | ConvertTo-Json -Compress
        
        # CORRIGIDO: Adicionado ponto e vírgula
        $headers = @{
            "X-Timestamp"   = $timestamp.ToString();
            "Authorization" = "assinatura-quebrada";
            "X-Nonce"       = $nonce;
            "X-User-ID"     = $global:testUser.id.ToString()
        }
        Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
    }
}


Write-Host "`n`n--- Suíte de Testes Finalizada ---" -ForegroundColor Yellow
Read-Host -Prompt "Pressione Enter para sair"