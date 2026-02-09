# VibeSQL Website

Landing page for https://vibesql.online

## Deployment via Azure DevOps Pipeline

### 1. Create Azure Static Web App (One-time)

**Azure Portal:**
1. Go to https://portal.azure.com
2. Create Static Web App
   - Name: `vibesql-online`
   - Region: West US 2
   - Plan: **Free**
   - Deployment: **Other**
3. After creation, go to **Settings > Configuration**
4. Copy the **Deployment token**

### 2. Configure DevOps Pipeline

**Add deployment token:**
1. Azure DevOps > Vibe SQL Microserver project
2. Pipelines > Library
3. Create variable group: `vibesql-website-vars`
4. Add variable: `deployment_token` (mark as secret)
5. Paste the token from step 1

**Create pipeline:**
1. Pipelines > New Pipeline
2. Select: Azure Repos Git
3. Repository: Vibe SQL Microserver
4. Existing YAML: `/azure-pipelines-website.yml`
5. Save and Run

### 3. Get Deployment URL

After pipeline completes:
- Static Web App URL: `https://vibesql-online.azurestaticapps.net`
- This is your CNAME target for DNS

### 4. Configure DNS

Point your domain to the Static Web App:
```
CNAME  vibesql.online  ->  vibesql-online.azurestaticapps.net
```

Then add custom domain in Azure Static Web App settings.

## Local Preview

```bash
start website\index.html
```

## Files

- `index.html` - Landing page (IdealVibe branded)
- `README.md` - This file
