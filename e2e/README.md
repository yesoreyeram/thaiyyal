# E2E Tests for Thaiyyal Workflow Engine

This directory contains end-to-end tests for the Thaiyyal workflow builder and execution system using Playwright.

## Overview

The E2E tests demonstrate and validate the complete user workflow:

1. Opening the workflow builder interface
2. Creating a workflow by adding nodes
3. Connecting nodes to form a data flow
4. Executing the workflow
5. Viewing execution results in the resizable panel
6. Handling errors gracefully

## Test Files

### workflow.spec.ts

Main test suite covering:

- **Workflow Creation Test**: Complete workflow from start to finish
  - Opens workflow builder
  - Adds Number nodes and Add operation
  - Attempts to connect nodes
  - Runs the workflow
  - Verifies execution panel appears
  - Checks results display
  - Tests panel resize functionality
  - Tests panel close functionality

- **Error Handling Test**: Validates error states
  - Runs empty/invalid workflow
  - Verifies error display
  - Captures error screenshots

## Screenshots

All test screenshots are saved to `e2e/screenshots/` directory:

| Screenshot | Description |
|------------|-------------|
| `01-workflow-builder.png` | Initial workflow builder interface |
| `02-node-palette.png` | Node palette opened with search |
| `03-nodes-on-canvas.png` | Nodes added to the canvas |
| `04-before-connections.png` | Workflow state before connections |
| `05-before-run.png` | Workflow ready to run |
| `06-execution-panel-loading.png` | Execution panel with loading indicator |
| `07-execution-results.png` | Execution results displayed |
| `08-final-state.png` | Final workflow state |
| `09-resized-panel.png` | Resized execution panel |
| `10-panel-closed.png` | State after closing panel |
| `11-error-handling.png` | Error state display |

## Prerequisites

1. **Build the application first**:
   ```bash
   ./build.sh
   ```

2. **Ensure server binary is ready**:
   ```bash
   ls -la ./server
   ```

3. **Install Playwright browsers** (first time only):
   ```bash
   npx playwright install chromium
   ```

## Running Tests

### Run all E2E tests
```bash
npm run test:e2e
```

### Run in UI mode (interactive)
```bash
npm run test:e2e:ui
```

### Run specific test file
```bash
npx playwright test e2e/workflow.spec.ts
```

### Run with headed browser (see the browser)
```bash
npx playwright test --headed
```

### Debug mode
```bash
npx playwright test --debug
```

## Test Configuration

Configuration is in `playwright.config.ts`:

- **Base URL**: http://localhost:8080
- **Test Directory**: `./e2e`
- **Browsers**: Chromium (desktop)
- **Screenshots**: Captured at each step
- **Video**: Recorded on failure
- **Web Server**: Automatically starts `./server -addr :8080`

## Viewing Results

### HTML Report
After running tests, view the HTML report:
```bash
npx playwright show-report e2e-report
```

### Screenshots
Screenshots are saved in `e2e/screenshots/` directory and can be viewed directly.

### Videos
Videos of failed tests are saved in `test-results/` directory.

## Test Features

### Automatic Server Start
The tests automatically start the Thaiyyal server before running and shut it down after completion. No manual server management needed.

### Screenshot Capture
Screenshots are taken at every major step:
- Before and after user interactions
- At state transitions
- On errors

### Error Handling
Tests gracefully handle:
- Missing elements (with appropriate timeouts)
- Network errors
- Execution failures
- UI state changes

### Resilient Selectors
Tests use multiple selector strategies:
- Text content
- ARIA roles
- CSS classes
- Fallback options

## Troubleshooting

### Server won't start
Ensure the server binary is built:
```bash
./build.sh
```

### Port already in use
Check if port 8080 is available:
```bash
lsof -i :8080
```

Kill any existing process or change the port in `playwright.config.ts`.

### Tests timeout
Increase timeout in individual tests or config:
```typescript
test('my test', async ({ page }) => {
  test.setTimeout(60000); // 60 seconds
  // ...
});
```

### Screenshots not captured
Ensure `e2e/screenshots/` directory exists:
```bash
mkdir -p e2e/screenshots
```

## Continuous Integration

For CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Install dependencies
  run: npm ci

- name: Build application
  run: ./build.sh

- name: Install Playwright browsers
  run: npx playwright install --with-deps chromium

- name: Run E2E tests
  run: npm run test:e2e

- name: Upload screenshots
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: e2e-screenshots
    path: e2e/screenshots/

- name: Upload test report
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: e2e-report
    path: e2e-report/
```

## Best Practices

1. **Wait for elements**: Always wait for elements before interacting
2. **Use test.step()**: Organize tests into logical steps
3. **Capture screenshots**: Take screenshots at key moments
4. **Handle race conditions**: Use appropriate timeouts and wait strategies
5. **Clean state**: Each test should start with a clean state
6. **Descriptive names**: Use clear, descriptive test names

## Future Enhancements

Potential test additions:

- [ ] Test node drag-and-drop positioning
- [ ] Test edge creation between nodes
- [ ] Test node configuration editing
- [ ] Test workflow save/load functionality
- [ ] Test workflow examples loading
- [ ] Test keyboard shortcuts
- [ ] Test different node types
- [ ] Test error scenarios
- [ ] Test performance with large workflows
- [ ] Test mobile/responsive views

## References

- [Playwright Documentation](https://playwright.dev/)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Thaiyyal Workflow Execution Guide](../docs/WORKFLOW_EXECUTION_GUIDE.md)
- [Thaiyyal API Documentation](../docs/API_EXAMPLES.md)
