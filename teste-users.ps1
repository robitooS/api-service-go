$baseUrl = "http://localhost:8080"
$hmacSecretBase64 = "v0ggs5DQqRPs7/sGSFKBhsaKZx5eb5eYVS3uYjZH+mU="
$hmacKeyBytes = [System.Convert]::FromBase64String($hmacSecretBase64)

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
                try {
                    $errorBody | ConvertFrom-Json | ConvertTo-Json -Depth 5 | Write-Host
                } catch { Write-Host $errorBody }
            } else {
                Write-Host "  [FALHA] Ocorreu um erro na requisição!" -ForegroundColor Red
                Write-Host "  Status Code: $statusCode (Esperado: $ExpectedStatusCode)" -ForegroundColor Yellow
                $errorResponseStream = $exception.Response.GetResponseStream()
                $streamReader = New-Object System.IO.StreamReader($errorResponseStream)
                $errorBody = $streamReader.ReadToEnd()
                $streamReader.Close()
                Write-Host "  Corpo da Resposta de Erro:"
                try {
                    $errorBody | ConvertFrom-Json | ConvertTo-Json -Depth 5 | Write-Host
                } catch { Write-Host $errorBody }
            }
        } else {
            Write-Host "  [ERRO CRÍTICO] Não foi possível conectar à API ou ocorreu um erro inesperado." -ForegroundColor Red
            Write-Host "  Detalhes: $($exception.Message)"
        }
    }
}

Write-Host "--- Iniciando Suíte de Testes da API ---" -ForegroundColor Yellow

$uniqueId = (Get-Date).Ticks
$global:testUser = @{
    name = "Usuario Teste $uniqueId"
    email = "teste$uniqueId@exemplo.com"
    password = "Password@123"
    id = 0
}

Write-Host "`n`n--- MÓDULO 1: CRIAÇÃO DE USUÁRIOS (/users/create) ---" -ForegroundColor Magenta

$createdUser = Run-Test -TestName "1.1 - Deve criar um usuário com sucesso" -ExpectedStatusCode 201 -TestAction {
    $body = @{ user_name = $global:testUser.name; user_email = $global:testUser.email; user_password = $global:testUser.password } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
}
if ($createdUser) {
    $global:testUser.id = $createdUser.ID
}

Run-Test -TestName "1.2 - Deve retornar erro ao tentar criar usuário com e-mail duplicado" -ExpectedStatusCode 500 -TestAction {
    $body = @{ user_name = "Outro Nome"; user_email = $global:testUser.email; user_password = "Password@456" } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
}

Run-Test -TestName "1.3 - Deve retornar erro para nome de usuário muito curto" -ExpectedStatusCode 500 -TestAction {
    $body = @{ user_name = "A"; user_email = "outroemail$uniqueId@exemplo.com"; user_password = "Password@123" } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
}

Run-Test -TestName "1.4 - Deve retornar erro para e-mail com formato inválido" -ExpectedStatusCode 500 -TestAction {
    $body = @{ user_name = "Nome Valido"; user_email = "email-invalido"; user_password = "Password@123" } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
}

Run-Test -TestName "1.5 - Deve retornar erro para senha fraca (sem número)" -ExpectedStatusCode 500 -TestAction {
    $body = @{ user_name = "Nome Valido"; user_email = "outroemail2-$uniqueId@exemplo.com"; user_password = "Password@" } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
}

Write-Host "`n`n--- MÓDULO 2: AUTENTICAÇÃO DE USUÁRIOS (/users/login) ---" -ForegroundColor Magenta

$loginResult = Run-Test -TestName "2.1 - Deve realizar login com sucesso" -ExpectedStatusCode 200 -TestAction {
    $body = @{ user_email = $global:testUser.email; user_password = $global:testUser.password } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/login" -Method Post -Body $body -ContentType "application/json"
}

Run-Test -TestName "2.2 - Deve retornar erro ao tentar login com e-mail inexistente" -ExpectedStatusCode 401 -TestAction {
    $body = @{ user_email = "naoexiste@exemplo.com"; user_password = "Password@123" } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/login" -Method Post -Body $body -ContentType "application/json"
}

Run-Test -TestName "2.3 - Deve retornar erro ao tentar login com senha incorreta" -ExpectedStatusCode 401 -TestAction {
    $body = @{ user_email = $global:testUser.email; user_password = "senhaerrada" } | ConvertTo-Json -Compress
    Invoke-RestMethod -Uri "$baseUrl/users/login" -Method Post -Body $body -ContentType "application/json"
}


Write-Host "`n`n--- MÓDULO 3: SEGURANÇA E ROTA PROTEGIDA (/users/get) ---" -ForegroundColor Magenta

$nonceValido = [System.Guid]::NewGuid().ToString()
Run-Test -TestName "3.1 - Deve aceder à rota protegida com autenticação HMAC válida" -ExpectedStatusCode 200 -TestAction {
    $method = "POST"
    $path = "/users/get"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $body = @{ user_id = $global:testUser.id } | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${body}:${nonceValido}"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes

    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = $signature
        "X-Nonce"       = $nonceValido
        "X-User-ID"     = $global:testUser.id.ToString()
    }

    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
}

Run-Test -TestName "3.2 - [SEGURANÇA] Deve bloquear ataque de replay (reutilização de nonce)" -ExpectedStatusCode 401 -TestAction {
    $method = "POST"
    $path = "/users/get"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $body = @{ user_id = $global:testUser.id } | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${body}:${nonceValido}"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes

    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = $signature
        "X-Nonce"       = $nonceValido
        "X-User-ID"     = $global:testUser.id.ToString()
    }

    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
}

Run-Test -TestName "3.3 - [SEGURANÇA] Deve bloquear acesso com assinatura HMAC inválida" -ExpectedStatusCode 401 -TestAction {
    $method = "POST"
    $path = "/users/get"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $nonce = [System.Guid]::NewGuid().ToString()
    $body = @{ user_id = $global:testUser.id } | ConvertTo-Json -Compress
    
    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = "assinatura-falsa"
        "X-Nonce"       = $nonce
        "X-User-ID"     = $global:testUser.id.ToString()
    }

    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
}

Run-Test -TestName "3.4 - [SEGURANÇA] Deve bloquear acesso com timestamp expirado" -ExpectedStatusCode 401 -TestAction {
    $method = "POST"
    $path = "/users/get"
    $timestamp = ([System.DateTimeOffset]::UtcNow).AddMinutes(-10).ToUnixTimeSeconds() 
    $nonce = [System.Guid]::NewGuid().ToString()
    $body = @{ user_id = $global:testUser.id } | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${body}:${nonce}"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes

    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = $signature
        "X-Nonce"       = $nonce
        "X-User-ID"     = $global:testUser.id.ToString()
    }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
}

Run-Test -TestName "3.5 - Deve retornar erro na ausência do header X-Nonce" -ExpectedStatusCode 400 -TestAction {
    $method = "POST"
    $path = "/users/get"
    $timestamp = [System.DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
    $body = @{ user_id = $global:testUser.id } | ConvertTo-Json -Compress
    $message = "${method}:${path}:${timestamp}:${body}:"
    $signature = Get-HmacSignature -Message $message -Key $hmacKeyBytes

    $headers = @{
        "X-Timestamp"   = $timestamp.ToString()
        "Authorization" = $signature
        "X-User-ID"     = $global:testUser.id.ToString()
    }
    Invoke-RestMethod -Uri "$baseUrl$path" -Method $method -Headers $headers -Body $body -ContentType "application/json"
}

Write-Host "`n`n--- MÓDULO 4: TESTE DE CARGA LEVE ---" -ForegroundColor Magenta
Write-Host "Iniciando criação de 10 usuários em paralelo..."

$jobs = @()
for ($i = 1; $i -le 10; $i++) {
    $job = Start-Job -ScriptBlock {
        param($baseUrl, $i)
        
        $jobUniqueId = (Get-Date).Ticks + $i 
        $email = "cargateste$jobUniqueId@exemplo.com"
        $name = "Usuario Carga $jobUniqueId"
        $password = "Password@123"
        $body = @{ user_name = $name; user_email = $email; user_password = $password } | ConvertTo-Json -Compress
        try {
            Invoke-RestMethod -Uri "$baseUrl/users/create" -Method Post -Body $body -ContentType "application/json"
            return "Usuário $i ($email) criado com sucesso."
        } catch {
            $statusCode = [int]$_.Exception.Response.StatusCode
            return "Falha ao criar usuário $i ($email). Status: $statusCode"
        }
    } -ArgumentList $baseUrl, $i
    $jobs += $job
}

$jobs | Wait-Job | ForEach-Object { Receive-Job $_ }
Write-Host "Teste de carga finalizado." -ForegroundColor Green


Write-Host "`n`n--- Suíte de Testes Finalizada ---" -ForegroundColor Yellow
Read-Host -Prompt "Pressione Enter para sair"