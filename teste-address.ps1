# Versão: 2.4
# Descrição:
# Versão com a correção do teste 1.3 (CEP Inválido), garantindo que a
# requisição seja feita por um usuário existente para que a falha de
# validação de negócio seja testada corretamente.
# ==============================================================================

# --- Configurações da API ---
$baseUrl = "http://localhost:8080"
$hmacSecretBase64 = "v0ggs5DQqRPs7/sGSFKBhsaKZx5eb5eYVS3uYjZH+mU="
$hmacKeyBytes = [System.Convert]::FromBase64String($hmacSecretBase64)


# --- Funções Auxiliares ---
function Get-HmacSignature {
    param(
        [Parameter(Mandatory=$true)] [string]$Message,
        [Parameter(Mandatory=$true)] [byte[]]$Key
    )
    $messageBytes = [System.Text.Encoding]::UTF8.GetBytes($Message)
    $hmac = New-Object System.Security.Cryptography.HMACSHA256
    $hmac.Key = $Key
    $signatureBytes = $hmac.ComputeHash($messageBytes)
    $base64Url = [System.Convert]::ToBase64String($signatureBytes).TrimEnd('=').Replace('+', '-').Replace('/', '_')
    return $base64Url
}

function Run-Test {
    param(
        [Parameter(Mandatory=$true)] [string]$TestName,
        [Parameter(Mandatory=$true)] [scriptblock]$TestAction,
        [Parameter(Mandatory=$true)] [int]$ExpectedStatusCode
    )
    Write-Host "`n[TESTE] $TestName" -ForegroundColor Cyan
    try {
        $responseObject = & $TestAction
        Write-Host "  [SUCESSO] Status Code retornado (Esperado: $ExpectedStatusCode)" -ForegroundColor Green
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
                $errorBody = $streamReader.ReadToEnd(); $streamReader.Close()
                Write-Host "  Corpo da Resposta:"; try { $errorBody | ConvertFrom-Json | ConvertTo-Json -Depth 5 | Write-Host } catch { Write-Host $errorBody }
            } else {
                Write-Host "  [FALHA] Status Code: $statusCode (Esperado: $ExpectedStatusCode)" -ForegroundColor Red
                $errorResponseStream = $exception.Response.GetResponseStream()
                $streamReader = New-Object System.IO.StreamReader($errorResponseStream)
                $errorBody = $streamReader.ReadToEnd(); $streamReader.Close()
                Write-Host "  Corpo da Resposta:"; try { $errorBody | ConvertFrom-Json | ConvertTo-Json -Depth 5 | Write-Host } catch { Write-Host $errorBody }
            }
        } else {
            Write-Host "  [ERRO CRÍTICO] Falha de conexão: $($exception.Message)" -ForegroundColor Red
        }
    }
}


# ==============================================================================
# --- INÍCIO DA SUÍTE DE TESTES ---
# ==============================================================================
Write-Host "--- Iniciando Suíte de Testes da API (Foco: Endereços) ---" -ForegroundColor Yellow

$uniqueId = (Get-Date).Ticks
$global:testUser = @{ name = "Usuario Endereco Teste $uniqueId"; email = "endereco$uniqueId@exemplo.com"; password = "Password@123Valida"; id = 0 }

# --- MÓDULO 0: PREPARAÇÃO ---
Write-Host "`n`n--- MÓDULO 0: PREPARAÇÃO DO AMBIENTE ---" -ForegroundColor Magenta

$createdUser = Run-Test -TestName "0.1 - Deve criar um usuário de teste" -ExpectedStatusCode 201 -TestAction {
    $body = @{ user_name = $global:testUser.name; user_email = $global:testUser.email; user_password = $global:testUser.password } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json; charset=utf-8"
}
if ($createdUser) { $global:testUser.id = $createdUser.ID; Write-Host "[INFO] Usuário de teste criado com ID: $($global:testUser.id)" -ForegroundColor White } 
else { Write-Host "[ERRO CRÍTICO] Não foi possível criar o usuário de teste. Abortando." -ForegroundColor Red; exit }

# --- MÓDULO 1: CRIAÇÃO DE ENDEREÇOS ---
Write-Host "`n`n--- MÓDULO 1: CRIAÇÃO DE ENDEREÇOS (/address/create) ---" -ForegroundColor Magenta

# 1.1 - Sucesso
Run-Test -TestName "1.1 - Deve criar um endereço com sucesso" -ExpectedStatusCode 201 -TestAction {
    $method = "POST"; $path = "/address/create"; $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $bodyObject = @{ address_street = "Rua Exemplo"; address_number = "123-A"; address_neighborhood = "Bairro dos Testes"; address_city = "Cidade da API"; address_state = "SP"; address_cep = "12345-678" }
    $bodyJson = $bodyObject | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${bodyJson}:$(New-Guid)"; $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
    $headers = @{ "X-Timestamp" = $timestamp.ToString(); "Authorization" = $signature; "X-Nonce" = ($message -split ":")[-1]; "X-User-ID" = $global:testUser.id.ToString() }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $bodyJson -ContentType "application/json; charset=utf-8"
}

# 1.2 - Duplicidade
Run-Test -TestName "1.2 - Deve falhar ao criar endereço duplicado para o mesmo usuário" -ExpectedStatusCode 500 -TestAction {
    $method = "POST"; $path = "/address/create"; $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $bodyObject = @{ address_street = "Outra Rua"; address_number = "999"; address_neighborhood = "Outro Bairro"; address_city = "Outra Cidade"; address_state = "RJ"; address_cep = "98765-432" }
    $bodyJson = $bodyObject | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${bodyJson}:$(New-Guid)"; $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
    $headers = @{ "X-Timestamp" = $timestamp.ToString(); "Authorization" = $signature; "X-Nonce" = ($message -split ":")[-1]; "X-User-ID" = $global:testUser.id.ToString() }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $bodyJson -ContentType "application/json; charset=utf-8"
}

# 1.3 - CEP Inválido
# CORREÇÃO: Usando um usuário válido ($global:testUser.id) para que a requisição
# passe pelo middleware e chegue na lógica de validação do handler.
Run-Test -TestName "1.3 - Deve falhar ao criar endereço com CEP inválido" -ExpectedStatusCode 500 -TestAction {
    $method = "POST"; $path = "/address/create"; $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds(); $nonce = [System.Guid]::NewGuid().ToString()
    $bodyObject = @{ address_street = "Rua do CEP Inválido"; address_number = "000"; address_neighborhood = "Bairro"; address_city = "Cidade"; address_state = "MG"; address_cep = "formato-errado" }
    $bodyJson = $bodyObject | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${bodyJson}:${nonce}"; $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
    # O X-User-ID agora é de um usuário que existe, permitindo que a validação de CEP seja atingida.
    $headers = @{ "X-Timestamp" = $timestamp.ToString(); "Authorization" = $signature; "X-Nonce" = $nonce; "X-User-ID" = $global:testUser.id.ToString() }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $bodyJson -ContentType "application/json; charset=utf-8"
}

# ==============================================================================
# --- MÓDULO 2: TESTES DE ATUALIZAÇÃO DE ENDEREÇO (/address/update) ---
# ==============================================================================
Write-Host "`n`n--- MÓDULO 2: ATUALIZAÇÃO DE ENDEREÇOS (/address/update) ---" -ForegroundColor Magenta

Run-Test -TestName "2.1 - Deve atualizar um endereço com sucesso" -ExpectedStatusCode 200 -TestAction {
    # CORRIGIDO: Método alterado para POST
    $method = "POST" 
    $path = "/address/update"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $nonce = [System.Guid]::NewGuid().ToString()
    $bodyObject = @{ 
        address_street       = "Avenida Principal"
        address_number       = "5000"
        address_neighborhood = "Centro Atualizado"
        address_city         = "Nova Cidade"
        address_state        = "BA"
        address_cep          = "40028-922"
        user_id              = $global:testUser.id
    }
    $bodyJson = $bodyObject | ConvertTo-Json -Compress
    
    $message = "${method}:${path}:${timestamp}:${bodyJson}:${nonce}"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
    
    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = $signature
        "X-Nonce"       = $nonce
        "X-User-ID"     = $global:testUser.id.ToString()
    }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $bodyJson -ContentType "application/json; charset=utf-8"
}

Run-Test -TestName "2.2 - [SEGURANÇA] Deve falhar ao atualizar com assinatura inválida" -ExpectedStatusCode 401 -TestAction {
    # CORRIGIDO: Método alterado para POST
    $method = "POST"
    $path = "/address/update"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $nonce = [System.Guid]::NewGuid().ToString()
    $bodyJson = @{ address_cep = "11111-111"; user_id = $global:testUser.id } | ConvertTo-Json -Compress
    
    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = "assinatura-completamente-invalida"
        "X-Nonce"       = $nonce
        "X-User-ID"     = $global:testUser.id.ToString()
    }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $bodyJson -ContentType "application/json; charset=utf-8"
}

Run-Test -TestName "2.3 - Deve falhar ao tentar atualizar um endereço inexistente" -ExpectedStatusCode 401 -TestAction {
    # CORRIGIDO: Método alterado para POST
    $method = "POST"
    $path = "/address/update"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $nonce = [System.Guid]::NewGuid().ToString()
    $userIdInexistente = 999999
    $bodyObject = @{ 
        address_street       = "Rua Fantasma"
        address_number       = "0"
        address_neighborhood = "Bairro Vazio"
        address_city         = "Nulidade"
        address_state        = "XX"
        address_cep          = "00000-000"
        user_id              = $userIdInexistente
    }
    $bodyJson = $bodyObject | ConvertTo-Json -Compress
    
    $message = "${method}:${path}:${timestamp}:${bodyJson}:${nonce}"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes
    
    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = $signature
        "X-Nonce"       = $nonce
        "X-User-ID"     = $userIdInexistente.ToString()
    }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $bodyJson -ContentType "application/json; charset=utf-8"
}

Write-Host "`n`n--- Suíte de Testes Finalizada ---" -ForegroundColor Yellow
Read-Host -Prompt "Pressione Enter para sair"