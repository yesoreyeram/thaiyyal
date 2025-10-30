# Frontend Pages

This directory contains the main pages for the Thaiyyal workflow builder application.

## Active Pages

### Home Page
- **`page.tsx`**: Main landing page
  - Dark theme with modern design
  - Hero section with feature highlights
  - Navigation to workflow builder
  - Workflow management features

### Workflow Builder
- **`workflow/page.tsx`**: Interactive workflow builder
  - Collapsible, floating node palette (hidden by default)
  - 7 categorized node groups with 23+ node types
  - Dark theme interface
  - Real-time JSON payload generation
  - Drag-and-drop workflow building

## Removed Pages

The following pages were removed as part of cleanup:
- `page-enhanced.tsx` - Legacy experimental version
- `page-original.tsx` - Original simplified version
- `tests/page.tsx` - Test scenarios page
- `pagination-tests/page.tsx` - Pagination test page

Legacy versions can be accessed via git history if needed.

## Page Structure

```
src/app/
├── page.tsx                    # Home page
├── workflow/
│   └── page.tsx               # Workflow builder
└── layout.tsx                 # Root layout
```
