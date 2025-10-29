# Frontend Page Variants

This directory contains multiple page variants for historical/comparison purposes:

## Active Page
- **`page.tsx`**: The main, canonical workflow builder page (23,622 LOC)
  - Full-featured implementation with all 23 node types
  - Context menu support for node operations
  - Right-click functionality
  - Comprehensive node palette organized by category

## Legacy/Experimental Pages (For Reference)
- **`page-original.tsx`**: Original simplified version (7,352 LOC)
  - Basic workflow builder with limited node types
  - Kept for comparison and rollback if needed
  
- **`page-enhanced.tsx`**: Enhanced experimental version (10,711 LOC)
  - Intermediate feature set
  - Testing ground for new features before merging to main

## Recommendation
Once the current implementation (page.tsx) is stable and verified, consider:
1. Archiving page-original.tsx and page-enhanced.tsx to a `legacy/` directory
2. Or removing them entirely if no longer needed
3. Use git history for reference to previous versions

## Migration Notes
- Current active: `page.tsx`
- All new features should be added to `page.tsx`
- Test page at: `/tests` (src/app/tests/page.tsx)
