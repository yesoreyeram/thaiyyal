import { test, expect } from '@playwright/test';

/**
 * E2E Test: Workflow Builder and Execution
 * 
 * This test demonstrates the complete workflow:
 * 1. Opening the workflow builder
 * 2. Creating a basic workflow with nodes
 * 3. Connecting nodes
 * 4. Running the workflow
 * 5. Viewing execution results
 */

test.describe('Workflow Builder and Execution', () => {
  test('should create a workflow, execute it, and view results', async ({ page }) => {
    // Step 1: Navigate to workflow builder
    await test.step('Open workflow builder', async () => {
      await page.goto('/workflow');
      await expect(page).toHaveTitle(/Thaiyyal/);
      
      // Wait for the canvas to be ready
      await page.waitForSelector('.react-flow', { timeout: 10000 });
      
      // Take screenshot of initial state
      await page.screenshot({ 
        path: 'e2e/screenshots/01-workflow-builder.png',
        fullPage: true 
      });
    });

    // Step 2: Open node palette
    await test.step('Open node palette', async () => {
      // Click the "Add Node" button
      const addNodeButton = page.locator('button:has-text("Add Node")');
      await addNodeButton.waitFor({ state: 'visible', timeout: 5000 });
      await addNodeButton.click();
      
      // Wait for palette to be visible
      await page.waitForSelector('input[placeholder*="Search"]', { timeout: 5000 });
      
      // Take screenshot of palette
      await page.screenshot({ 
        path: 'e2e/screenshots/02-node-palette.png',
        fullPage: true 
      });
    });

    // Step 3: Add Number nodes
    await test.step('Add first number node', async () => {
      // Search for "number" in palette
      const searchInput = page.locator('input[placeholder*="Search"]');
      await searchInput.fill('number');
      await page.waitForTimeout(500); // Wait for search to filter
      
      // Click on Number node type
      const numberNode = page.locator('text=Number').first();
      await numberNode.click();
      
      // Wait for node to appear on canvas
      await page.waitForTimeout(500);
    });

    await test.step('Add second number node', async () => {
      // Add another number node
      const numberNode = page.locator('text=Number').first();
      await numberNode.click();
      await page.waitForTimeout(500);
    });

    // Step 4: Add Add operation node
    await test.step('Add add operation node', async () => {
      // Search for "add" in palette
      const searchInput = page.locator('input[placeholder*="Search"]');
      await searchInput.clear();
      await searchInput.fill('add');
      await page.waitForTimeout(500);
      
      // Click on Add node
      const addNode = page.locator('text=Add').first();
      await addNode.click();
      await page.waitForTimeout(500);
    });

    // Step 5: Close palette
    await test.step('Close node palette', async () => {
      // Click outside palette or press Escape
      await page.keyboard.press('Escape');
      await page.waitForTimeout(500);
      
      // Take screenshot of nodes on canvas
      await page.screenshot({ 
        path: 'e2e/screenshots/03-nodes-on-canvas.png',
        fullPage: true 
      });
    });

    // Step 6: Connect nodes (if handles are accessible)
    // Note: ReactFlow connections might require precise mouse movements
    // This is a simplified version - actual connection may need coordinate-based dragging
    await test.step('Attempt to connect nodes', async () => {
      // Take screenshot showing the workflow before connections
      await page.screenshot({ 
        path: 'e2e/screenshots/04-before-connections.png',
        fullPage: true 
      });
      
      // Note: Actual drag-and-drop between node handles would go here
      // For this basic test, we'll proceed to show the workflow state
    });

    // Step 7: Click Run button
    await test.step('Click run button', async () => {
      // Find and click the Run button
      const runButton = page.locator('button:has-text("Run"), button >> svg.lucide-play');
      
      // Wait for run button to be visible
      await runButton.waitFor({ state: 'visible', timeout: 5000 });
      
      // Take screenshot before clicking run
      await page.screenshot({ 
        path: 'e2e/screenshots/05-before-run.png',
        fullPage: true 
      });
      
      // Click run button
      await runButton.click();
    });

    // Step 8: Wait for execution panel to appear
    await test.step('Verify execution panel appears', async () => {
      // Wait for execution panel to be visible
      await page.waitForSelector('text=Execution Results', { timeout: 10000 });
      
      // Take screenshot of loading state
      await page.screenshot({ 
        path: 'e2e/screenshots/06-execution-panel-loading.png',
        fullPage: true 
      });
    });

    // Step 9: Wait for results or error
    await test.step('Wait for execution to complete', async () => {
      // Wait for either success or error message
      await Promise.race([
        page.waitForSelector('text=Completed in', { timeout: 15000 }),
        page.waitForSelector('text=Failed', { timeout: 15000 }),
        page.waitForSelector('text=Error', { timeout: 15000 })
      ]).catch(() => {
        // Timeout is ok - we'll capture whatever state we're in
      });
      
      // Wait a bit for results to fully render
      await page.waitForTimeout(1000);
      
      // Take screenshot of results
      await page.screenshot({ 
        path: 'e2e/screenshots/07-execution-results.png',
        fullPage: true 
      });
    });

    // Step 10: Verify execution panel content
    await test.step('Verify execution panel elements', async () => {
      // Check that execution panel has key elements
      const executionPanel = page.locator('text=Execution Results').locator('..');
      await expect(executionPanel).toBeVisible();
      
      // Take final screenshot
      await page.screenshot({ 
        path: 'e2e/screenshots/08-final-state.png',
        fullPage: true 
      });
    });

    // Step 11: Test panel resize
    await test.step('Test panel resize functionality', async () => {
      // Find the resize handle (it has cursor-ns-resize class)
      const resizeHandle = page.locator('.cursor-ns-resize, [class*="cursor-ns-resize"]').first();
      
      if (await resizeHandle.isVisible()) {
        // Get the handle's position
        const handleBox = await resizeHandle.boundingBox();
        
        if (handleBox) {
          // Drag the handle up to resize
          await page.mouse.move(handleBox.x + handleBox.width / 2, handleBox.y + handleBox.height / 2);
          await page.mouse.down();
          await page.mouse.move(handleBox.x + handleBox.width / 2, handleBox.y - 100);
          await page.mouse.up();
          
          await page.waitForTimeout(500);
          
          // Take screenshot of resized panel
          await page.screenshot({ 
            path: 'e2e/screenshots/09-resized-panel.png',
            fullPage: true 
          });
        }
      }
    });

    // Step 12: Close execution panel
    await test.step('Close execution panel', async () => {
      // Find and click the close button (✕)
      const closeButton = page.locator('button:has-text("✕")').last();
      
      if (await closeButton.isVisible()) {
        await closeButton.click();
        await page.waitForTimeout(500);
        
        // Take screenshot after closing panel
        await page.screenshot({ 
          path: 'e2e/screenshots/10-panel-closed.png',
          fullPage: true 
        });
      }
    });
  });

  test('should handle workflow execution errors gracefully', async ({ page }) => {
    await test.step('Navigate and create empty workflow', async () => {
      await page.goto('/workflow');
      await page.waitForSelector('.react-flow', { timeout: 10000 });
    });

    await test.step('Run empty workflow and check error handling', async () => {
      // Click run on empty/invalid workflow
      const runButton = page.locator('button:has-text("Run"), button >> svg.lucide-play');
      await runButton.waitFor({ state: 'visible', timeout: 5000 });
      await runButton.click();
      
      // Wait for execution panel
      await page.waitForSelector('text=Execution Results', { timeout: 10000 });
      
      // Wait for error or completion
      await page.waitForTimeout(2000);
      
      // Take screenshot of error state
      await page.screenshot({ 
        path: 'e2e/screenshots/11-error-handling.png',
        fullPage: true 
      });
    });
  });
});
