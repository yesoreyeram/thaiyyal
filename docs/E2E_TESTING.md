# E2E Testing with Playwright

This document explains the end-to-end testing setup for the Thaiyyal workflow engine using Playwright.

## Overview

The E2E tests provide automated testing and screenshot capture for the complete workflow creation and execution flow. These tests serve dual purposes:

1. **Automated Testing**: Validate the workflow builder and execution functionality
2. **Documentation**: Generate screenshots for the documentation automatically

## Architecture

```
thaiyyal/
├── e2e/                          # E2E test directory
│   ├── workflow.spec.ts          # Main workflow tests
│   ├── screenshots/              # Generated screenshots (committed)
│   │   ├── 01-workflow-builder.png
│   │   ├── 02-node-palette.png
│   │   ├── ... (11 total screenshots)
│   └── README.md                 # E2E documentation
├── playwright.config.ts          # Playwright configuration
├── test-results/                 # Test results (gitignored)
├── e2e-report/                   # HTML report (gitignored)
└── package.json                  # NPM scripts for tests
```

## Test Coverage

### 1. Workflow Creation and Execution Test

This comprehensive test covers the entire user journey:

**Steps:**

1. **Open Workflow Builder** (`01-workflow-builder.png`)
   - Navigate to `/workflow`
   - Verify page loads
   - Wait for ReactFlow canvas to initialize

2. **Open Node Palette** (`02-node-palette.png`)
   - Click "Add Node" button
   - Verify palette appears with search box

3. **Add Nodes** (`03-nodes-on-canvas.png`)
   - Add two Number nodes
   - Add one Add operation node
   - Close palette

4. **Prepare Workflow** (`04-before-connections.png`, `05-before-run.png`)
   - Show nodes on canvas
   - Prepare for execution

5. **Execute Workflow** (`06-execution-panel-loading.png`)
   - Click Run button
   - Verify execution panel opens
   - Show loading indicator

6. **View Results** (`07-execution-results.png`, `08-final-state.png`)
   - Wait for execution to complete
   - Display execution summary
   - Show final output and node results

7. **Test Panel Interactions** (`09-resized-panel.png`, `10-panel-closed.png`)
   - Resize panel by dragging
   - Close panel

### 2. Error Handling Test

Tests error scenarios:

**Steps:**

1. Navigate to workflow builder
2. Run empty/invalid workflow
3. Verify error display (`11-error-handling.png`)

## Setup

### Prerequisites

1. **Node.js** 20.x or later
2. **Built application**:
   ```bash
   ./build.sh
   ```

### Installation

1. **Install Playwright** (already in devDependencies):
   ```bash
   npm install
   ```

2. **Install browser binaries**:
   ```bash
   npx playwright install chromium
   ```

## Running Tests

### Basic Commands

```bash
# Run all E2E tests
npm run test:e2e

# Run with UI mode (recommended for development)
npm run test:e2e:ui

# Run with headed browser (see the browser)
npm run test:e2e:headed

# Debug mode (step through tests)
npm run test:e2e:debug

# View test report
npm run test:e2e:report
```

### Advanced Usage

```bash
# Run specific test file
npx playwright test e2e/workflow.spec.ts

# Run specific test by name
npx playwright test -g "should create a workflow"

# Run in specific browser
npx playwright test --project=chromium

# Update snapshots
npx playwright test --update-snapshots

# Run with trace
npx playwright test --trace on
```

## Configuration

### playwright.config.ts

Key configuration options:

```typescript
{
  testDir: './e2e',                    // Test directory
  baseURL: 'http://localhost:8080',    // Application URL
  fullyParallel: true,                 // Parallel execution
  retries: process.env.CI ? 2 : 0,    // Retry on CI
  reporter: [
    ['html', { outputFolder: 'e2e-report' }],
    ['list']
  ],
  use: {
    screenshot: 'only-on-failure',     // Screenshot on failure
    video: 'retain-on-failure',        // Video on failure
    trace: 'on-first-retry',          // Trace on retry
  },
  webServer: {
    command: './server -addr :8080',   // Auto-start server
    url: 'http://localhost:8080',
    timeout: 120 * 1000,
  },
}
```

## Screenshots

All screenshots are automatically captured and saved to `e2e/screenshots/`:

| File | Description | Usage |
|------|-------------|-------|
| `01-workflow-builder.png` | Initial workflow builder interface | Documentation, visual regression |
| `02-node-palette.png` | Node palette with search | Feature demonstration |
| `03-nodes-on-canvas.png` | Nodes added to canvas | Workflow creation guide |
| `04-before-connections.png` | Pre-connection state | Connection tutorial |
| `05-before-run.png` | Ready to execute | Execution guide |
| `06-execution-panel-loading.png` | Loading state | Loading indicator demo |
| `07-execution-results.png` | Execution results | Results display guide |
| `08-final-state.png` | Complete workflow | Final state reference |
| `09-resized-panel.png` | Resized execution panel | Panel resize feature |
| `10-panel-closed.png` | Panel closed state | Panel management |
| `11-error-handling.png` | Error state | Error handling guide |

### Using Screenshots in Documentation

Screenshots can be referenced in markdown:

```markdown
![Workflow Builder](../e2e/screenshots/01-workflow-builder.png)
```

Or copied to docs directory:

```bash
cp e2e/screenshots/*.png docs/screenshots/
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.7'
          
      - name: Install dependencies
        run: npm ci
        
      - name: Build application
        run: ./build.sh
        
      - name: Install Playwright Browsers
        run: npx playwright install --with-deps chromium
        
      - name: Run E2E tests
        run: npm run test:e2e
        
      - name: Upload screenshots
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: e2e-screenshots
          path: e2e/screenshots/
          retention-days: 30
          
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: e2e-test-results
          path: test-results/
          retention-days: 7
          
      - name: Upload HTML report
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: e2e-html-report
          path: e2e-report/
          retention-days: 7
```

## Troubleshooting

### Common Issues

#### 1. Server won't start

**Problem**: `Error: Command failed: ./server -addr :8080`

**Solutions**:
```bash
# Ensure server is built
./build.sh

# Check server binary exists
ls -la ./server

# Test server manually
./server -addr :8080
```

#### 2. Port already in use

**Problem**: `EADDRINUSE: address already in use :::8080`

**Solutions**:
```bash
# Find process using port
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port in playwright.config.ts
```

#### 3. Tests timeout

**Problem**: Tests fail with timeout errors

**Solutions**:
- Increase timeout in test:
  ```typescript
  test.setTimeout(60000); // 60 seconds
  ```
- Increase global timeout in config
- Check server is responding: `curl http://localhost:8080`

#### 4. Screenshots not captured

**Problem**: Screenshot files missing

**Solutions**:
```bash
# Ensure directory exists
mkdir -p e2e/screenshots

# Check permissions
chmod 755 e2e/screenshots

# Run with --headed to see what's happening
npm run test:e2e:headed
```

#### 5. Browser not installed

**Problem**: `browserType.launch: Executable doesn't exist`

**Solutions**:
```bash
# Install browsers
npx playwright install chromium

# Or install all browsers with system dependencies
npx playwright install --with-deps
```

## Best Practices

### 1. Wait Strategies

Use appropriate wait strategies:

```typescript
// Wait for element
await page.waitForSelector('.react-flow');

// Wait for navigation
await page.waitForNavigation();

// Wait for timeout (last resort)
await page.waitForTimeout(1000);

// Wait for function
await page.waitForFunction(() => window.loaded === true);
```

### 2. Selectors

Use resilient selectors:

```typescript
// Good - semantic selectors
page.locator('button:has-text("Run")')
page.locator('[aria-label="Add Node"]')
page.locator('role=button[name="Execute"]')

// Avoid - brittle selectors
page.locator('.button-class-123')
page.locator('div > div > button:nth-child(3)')
```

### 3. Screenshots

Take screenshots at key moments:

```typescript
// Full page screenshot
await page.screenshot({ 
  path: 'screenshot.png',
  fullPage: true 
});

// Element screenshot
await element.screenshot({ path: 'element.png' });

// With custom viewport
await page.setViewportSize({ width: 1920, height: 1080 });
await page.screenshot({ path: 'hd.png' });
```

### 4. Test Isolation

Each test should be independent:

```typescript
test.beforeEach(async ({ page }) => {
  // Reset state before each test
  await page.goto('/workflow');
});

test.afterEach(async ({ page }) => {
  // Cleanup after each test
  await page.close();
});
```

## Advanced Features

### Visual Regression Testing

Add visual comparison tests:

```typescript
test('visual regression', async ({ page }) => {
  await page.goto('/workflow');
  await expect(page).toHaveScreenshot('workflow-baseline.png');
});
```

### Custom Fixtures

Create reusable fixtures:

```typescript
import { test as base } from '@playwright/test';

const test = base.extend({
  workflowPage: async ({ page }, use) => {
    await page.goto('/workflow');
    await page.waitForSelector('.react-flow');
    await use(page);
  },
});
```

### Network Mocking

Mock API responses:

```typescript
await page.route('**/api/v1/workflow/execute', route => {
  route.fulfill({
    status: 200,
    body: JSON.stringify({
      success: true,
      results: { /* mock data */ }
    })
  });
});
```

## Maintenance

### Updating Tests

When UI changes:

1. Update selectors in tests
2. Re-run tests to generate new screenshots
3. Review and commit new screenshots
4. Update documentation if needed

### Updating Screenshots

To refresh all screenshots:

```bash
# Delete old screenshots
rm e2e/screenshots/*.png

# Run tests to regenerate
npm run test:e2e

# Review and commit
git add e2e/screenshots/*.png
git commit -m "Update E2E screenshots"
```

## Resources

- [Playwright Documentation](https://playwright.dev/)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Thaiyyal Documentation](../docs/)

## Contributing

When adding new E2E tests:

1. Follow existing test structure
2. Use descriptive test names
3. Add screenshots at key steps
4. Update this documentation
5. Ensure tests pass in CI

## Support

For issues with E2E tests:

1. Check [Troubleshooting](#troubleshooting)
2. Review [Playwright Docs](https://playwright.dev/)
3. Open an issue on GitHub
4. Check test logs: `cat test-results/*/test-log.txt`
