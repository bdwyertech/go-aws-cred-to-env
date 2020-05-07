# AWS Credential to Environment
Fetches AWS Credentials and loads them into the Shell.

```bash
. <(aws-cred-to-env)
```

```powershell
. $([System.Management.Automation.ScriptBlock]::Create($(aws-cred-to-env | Out-String)))
```
