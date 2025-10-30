# GitHub Pages Setup Instructions

This repository is now configured with GitHub Actions to automatically build and deploy the Thaiyyal frontend to GitHub Pages. Follow these steps to enable GitHub Pages for this repository.

## Prerequisites

- You must have admin access to this repository
- GitHub Pages must be enabled for your GitHub account/organization

## Step-by-Step Setup

### 1. Navigate to Repository Settings

1. Go to your repository on GitHub: `https://github.com/yesoreyeram/thaiyyal`
2. Click on the **Settings** tab (⚙️ icon at the top of the repository)

### 2. Enable GitHub Pages

1. In the left sidebar, scroll down and click on **Pages** (under "Code and automation")
2. Under **"Build and deployment"** section:
   - **Source**: Select **"GitHub Actions"** from the dropdown
   - This allows GitHub Actions workflows to deploy to Pages
3. Click **Save** (if required)

### 3. Verify Workflow Permissions (Important!)

1. In the left sidebar, click on **Actions** → **General**
2. Scroll down to **"Workflow permissions"**
3. Ensure the following settings are configured:
   - Select **"Read and write permissions"**
   - Check ✅ **"Allow GitHub Actions to create and approve pull requests"**
4. Click **Save** at the bottom

### 4. Trigger the Deployment

The workflow is configured to run automatically on:
- Every push to the `main` branch
- Manual trigger via workflow_dispatch

#### Option A: Merge this Pull Request

1. Merge this pull request into the `main` branch
2. The workflow will automatically run and deploy

#### Option B: Manual Trigger

1. Go to **Actions** tab in your repository
2. Click on **"Deploy to GitHub Pages"** workflow in the left sidebar
3. Click the **"Run workflow"** button
4. Select the `main` branch
5. Click **"Run workflow"**

### 5. Monitor the Deployment

1. Go to the **Actions** tab
2. You should see a new workflow run for "Deploy to GitHub Pages"
3. Click on the workflow run to see the progress
4. The deployment typically takes 1-3 minutes

### 6. Access Your Deployed Site

Once the deployment is complete:

1. Go to **Settings** → **Pages**
2. You'll see a message: **"Your site is live at https://yesoreyeram.github.io/thaiyyal/"**
3. Click the URL to visit your deployed Thaiyyal application

**Note**: If you're deploying to a project page (not a user/organization page), the URL will be:
- `https://<username>.github.io/<repository-name>/`
- Example: `https://yesoreyeram.github.io/thaiyyal/`

## Workflow File Location

The GitHub Actions workflow is located at:
```
.github/workflows/deploy-pages.yml
```

## Workflow Configuration

The workflow:
- ✅ Builds the Next.js application with static export
- ✅ Uploads the built files as an artifact
- ✅ Deploys to GitHub Pages
- ✅ Runs on every push to `main` branch
- ✅ Can be manually triggered

## Customization

### Change Deployment Branch

To deploy from a different branch, edit `.github/workflows/deploy-pages.yml`:

```yaml
on:
  push:
    branches:
      - main  # Change this to your desired branch
```

### Add Base Path

If your repository name is not the domain root, you may need to configure a base path in `next.config.ts`:

```typescript
const nextConfig: NextConfig = {
  devIndicators: false,
  output: 'export',
  basePath: '/thaiyyal',  // Add this line with your repo name
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
};
```

## Troubleshooting

### Deployment Failed

1. **Check workflow logs**: Go to Actions tab and click on the failed workflow run
2. **Common issues**:
   - Insufficient permissions: Verify step 3 above
   - Build errors: Check the build logs in the workflow
   - GitHub Pages not enabled: Verify step 2 above

### 404 Error on Deployed Site

1. **Check base path**: Ensure `next.config.ts` has the correct `basePath` if deploying to a project page
2. **Clear browser cache**: Force refresh with Ctrl+Shift+R (Windows/Linux) or Cmd+Shift+R (Mac)
3. **Wait a few minutes**: DNS propagation can take time

### Changes Not Reflecting

1. **Check workflow status**: Ensure the latest workflow run completed successfully
2. **Clear GitHub Pages cache**: Settings → Pages → Uncheck and re-check the source
3. **Hard refresh browser**: Ctrl+Shift+R or Cmd+Shift+R

## Additional Configuration

### Custom Domain

To use a custom domain:

1. Go to **Settings** → **Pages**
2. Under **"Custom domain"**, enter your domain name
3. Configure your DNS provider to point to GitHub Pages
4. Wait for DNS propagation (can take up to 24-48 hours)

### HTTPS

GitHub Pages automatically provides HTTPS for `*.github.io` domains. For custom domains:

1. Ensure DNS is configured correctly
2. Check the **"Enforce HTTPS"** checkbox in Pages settings
3. Wait for the certificate to be provisioned

## Maintenance

### Automatic Deployments

- Every push to `main` will automatically trigger a new deployment
- No manual intervention required

### Manual Deployments

- Use the "Run workflow" button in the Actions tab
- Useful for deploying without new commits

## Support

If you encounter issues:

1. Check the [GitHub Pages documentation](https://docs.github.com/en/pages)
2. Review the [GitHub Actions logs](https://github.com/yesoreyeram/thaiyyal/actions)
3. Open an issue in this repository

## Summary Checklist

- [ ] Navigate to Settings → Pages
- [ ] Set Source to "GitHub Actions"
- [ ] Navigate to Settings → Actions → General
- [ ] Set Workflow permissions to "Read and write"
- [ ] Merge this PR or manually trigger the workflow
- [ ] Wait for deployment to complete (1-3 minutes)
- [ ] Visit your site at `https://yesoreyeram.github.io/thaiyyal/`
- [ ] Verify the application loads correctly

---

**Congratulations!** 🎉 Your Thaiyyal application is now deployed to GitHub Pages!
