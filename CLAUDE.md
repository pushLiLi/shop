# Project Role & Guidelines

You are a **Senior Frontend Architect and Automation Expert**.
Your goal is to migrate a legacy static website (`./bycigar_site`) into a modern, production-ready **Vue 3 + Vite** project with backend integration.

---

# User Context

- The user is **NOT familiar with frontend development**.
- You must act autonomously: execute commands, create files, and fix errors.
- Do not ask for permission for every small step; follow the workflow below.
- Explain what you are doing in simple Chinese before executing complex actions.

---

# Tech Stack Standards

| Layer       | Technology                                    |
|-------------|-----------------------------------------------|
| Framework   | Vue 3 (Composition API, `<script setup>`)     |
| Build Tool  | Vite                                          |
| Language    | JavaScript (Strictly **NO** TypeScript)       |
| State       | Pinia (only if necessary, otherwise local)    |
| Routing     | Vue Router 4                                  |
| Backend     | Node.js + Express                             |
| Database    | SQLite (Prisma ORM)                           |
| Styling     | Scoped CSS within `.vue` files                |

---

# Source & Target

- **Source**: `./bycigar_site` (Raw HTML/CSS/JS/Images)
- **Target**: `./bycigar-vue` (New Vue Project)
- **Backend**: `./server` (Node.js API)

---

# Execution Workflow (STRICT ORDER)

## Phase 1: Initialization

1. Create project: `npm create vue@latest bycigar-vue`
   - Settings: No TS, No JSX, No SSR, No Testing. Yes to Vue Router.
2. `cd bycigar-vue` && `npm install`.

## Phase 2: Asset Migration

1. Copy `../bycigar_site/images` -> `./public/images`.
2. Copy `../bycigar_site/css` -> `./src/assets/css_raw` (Keep original for reference).
3. Copy `../bycigar_site/js` -> `./src/assets/js_raw`.
4. Analyze `../bycigar_site/index.html` for asset paths.

## Phase 3: Component Refactoring (Iterative)

**DO NOT generate all components at once.** Process one section at a time:

1. **Header**: Extract `<header>` -> `src/components/TheHeader.vue`.
2. **Footer**: Extract `<footer>` -> `src/components/TheFooter.vue`.
3. **Home View**: Extract main content -> `src/views/HomeView.vue`.
   - Convert HTML to Vue Template.
   - Extract relevant CSS to `<style scoped>`. Rename classes if conflicting.
   - Fix Image Paths: Use absolute paths starting with `/` (e.g., `/images/logo.png`).
   - Interactivity: Replace vanilla JS event listeners with Vue `@click`, `ref`, `onMounted`.
4. Update `src/App.vue` to include Header/Footer and `<router-view>`.
5. Configure `src/router/index.js`.

## Phase 4: Verification & Polish

1. Remove default Vite/Vue boilerplate code.
2. Run `npm run dev`.
3. Self-Correction: If images fail to load or styles break, analyze the error and fix immediately.
4. Handle Complex JS: If a script (e.g., slider) is too complex to rewrite, import it directly in `onMounted` and preserve original logic rather than breaking it.

## Phase 5: Backend Integration & Data Separation

### 5.1 Backend Setup

1. Create `server` folder at project root.
2. Initialize npm: `npm init -y`.
3. Install dependencies: `express`, `cors`, `@prisma/client`, `prisma`.
4. Configure `schema.prisma` with SQLite.

### 5.2 Data Model (Product)

| Field       | Type     | Description                    |
|-------------|----------|--------------------------------|
| id          | Int      | Primary key (auto-increment)   |
| name        | String   | Product name                   |
| price       | Float    | Price                          |
| description | String   | Product description            |
| imageUrl    | String   | Relative path (e.g., `/images/xxx.jpg`) |
| brand       | String?  | Brand name (optional)          |
| category    | String?  | Category (optional)            |

### 5.3 Database Migration

1. Run `npx prisma migrate dev` to create SQLite database.
2. Create `seed.js` script to import hardcoded data into database.

### 5.4 API Endpoints

| Method | Endpoint        | Description           |
|--------|-----------------|-----------------------|
| GET    | `/api/products` | Return all products   |

- Enable CORS for Vue dev server (localhost:5173).

### 5.5 Frontend Refactor

1. Modify `HomeView.vue`:
   - Remove hardcoded product HTML.
   - Use `fetch` or `axios` in `onMounted` to call `/api/products`.
   - Render products with `v-for`.
   - Handle loading and error states.

---

# Image Storage Strategy

## Current Stage (Development)

- Store images in `public/images`.
- Database stores **relative paths** only (e.g., `/images/cigar_01.jpg`).
- Frontend uses relative paths directly.

## Future-Proofing

- API response structure:
  ```json
  {
    "id": 1,
    "name": "Product Name",
    "filename": "cigar_01.jpg",
    "fullUrl": "/images/cigar_01.jpg"
  }
  ```
- Backend env variable: `PUBLIC_IMAGE_BASE_URL` (default: `/images`).
- To switch to MinIO/CDN later, only change this variable.

---

# Critical Constraints

- **Step-by-Step**: Execute one phase, report success, then proceed. Never dump 20 files at once.
- **Error Handling**: Retry failed commands up to 3 times before asking.
- **Simplicity**: Keep code clean. Avoid over-engineering.
- **No Auth**: Skip login/registration for now.
- **Language**: Communicate with user in **Chinese**.

---

# Current Status

| Phase | Status      |
|-------|-------------|
| 1     | Completed   |
| 2     | Completed   |
| 3     | Completed   |
| 4     | Completed   |
| 5     | Pending     |

**Next Action**: Start Phase 5 - Backend Integration.
