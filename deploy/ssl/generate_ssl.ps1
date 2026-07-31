# Self-Signed SSL Certificate Generator for IP 172.17.10.109
$IP_ADDRESS = "172.17.10.109"
$OUT_DIR = $PSScriptRoot

Write-Host "Membuat sertifikat SSL untuk IP: $IP_ADDRESS..."

$cert = New-SelfSignedCertificate -DnsName "localhost", $IP_ADDRESS -CertStoreLocation "cert:\CurrentUser\My" -NotAfter (Get-Date).AddYears(5) -KeyExportPolicy Exportable -KeySpec Signature -KeyLength 2048 -HashAlgorithm SHA256 -FriendlyName "Sistem AssetKu SSL"

$crtPath = Join-Path $OUT_DIR "server.crt"
$certBytes = $cert.Export([System.Security.Cryptography.X509Certificates.X509ContentType]::Cert)

$pemHeader = "-----BEGIN CERTIFICATE-----"
$pemFooter = "-----END CERTIFICATE-----"
$pemBody = [System.Convert]::ToBase64String($certBytes, [System.Base64FormattingOptions]::InsertLineBreaks)
$crtPemContent = "$pemHeader`n$pemBody`n$pemFooter"
[System.IO.File]::WriteAllText($crtPath, $crtPemContent)

Write-Host "✓ Public Key berhasil dibuat: $crtPath"

$pfxPath = Join-Path $OUT_DIR "server.pfx"
$pfxPassword = ConvertTo-SecureString -String "assetku123" -Force -AsPlainText
$cert | Export-PfxCertificate -FilePath $pfxPath -Password $pfxPassword | Out-Null

Write-Host "✓ PFX Key berhasil dibuat: $pfxPath"
