# SkuliApp — Master Plan

> A school management platform. This document is the single source of truth for the
> tech stack, architecture, feature scope, and build roadmap.

## Table of Contents

- [1. Overview](#1-overview)
- [2. Tech Stack](#2-tech-stack)
- [3. Build Roadmap](#3-build-roadmap)
- [4. Core Features](#4-core-features)
- [5. Project Structure](#5-project-structure)
  - [5.1 Frontend (React)](#51-frontend-react)
  - [5.2 Backend (Go)](#52-backend-go)
  - [5.3 Module Structure](#53-module-structure)
  - [5.4 Database Layer](#54-database-layer)
- [6. Roles Reference](#6-roles-reference)

---

## 1. Overview

SkuliApp is a multi-tenant school management system. Each school operates in an
isolated context, with role-based access controlling what every user can see and do.
The MVP focuses on one wedge feature taken end-to-end before broadening scope.

**Guiding principles**

- Ship a single wedge feature end-to-end before expanding.
- Keep every backend module self-contained (handler → service → repository).
- Deploy early and often — Go's single-binary output makes this cheap.

---

## 2. Tech Stack

| Layer        | Technology | Notes |
| ------------ | ---------- | ----- |
| **Frontend** | React + TypeScript | UI framework |
| | Vite | Build tool / dev server |
| | Tailwind CSS | Styling |
| **Backend**  | Go (Golang) | API services |
| | Gin or Fiber | Web framework |
| | GORM | ORM |
| | JWT | Authentication |
| | golang-migrate | Database migrations |
| **Database** | PostgreSQL | Primary datastore |
| **Storage**  | Local (dev) · AWS S3 (prod) | File uploads |
| **DevOps**   | Docker · GitHub Actions | Containers + CI/CD |
| **Hosting**  | Railway / Render (backend) · Vercel (frontend) | MVP hosting |

---

## 3. Build Roadmap

Suggested order of execution, from foundation to first shipped feature.

- [ ] **Phase 1 — Project scaffold:** Go module setup, Gin router, project structure (`/handlers`, `/models`, `/middleware`, `/db`).
- [ ] **Phase 2 — Data layer:** Postgres connection + GORM models — `School`, `User` (with role field), `Class`, `Student`, `Teacher`.
- [ ] **Phase 3 — Auth:** JWT issuing on login; middleware to verify tokens and attach role/school context to each request.
- [ ] **Phase 4 — Wedge API:** Core REST endpoints for the first wedge feature (CRUD + business logic).
- [ ] **Phase 5 — Frontend shell:** React app via Vite; auth flow (login, protected routes, role-based UI).
- [ ] **Phase 6 — Wire the wedge:** Connect the wedge feature end-to-end (e.g., fee records, attendance marking).
- [ ] **Phase 7 — Deploy:** Ship early — Go's single-binary deploy keeps this painless.

---

## 4. Core Features

### Authentication
- Login
- Logout
- Forgot password
- Reset password
- Change password
- JWT authentication

### User Roles
- Super Admin
- School Admin
- Principal
- Teacher
- Accountant
- Parent
- Student

### Student Management
- Student admission
- Student profiles
- Student transfers
- Student promotion
- Student documents
- Attendance

### Teacher Management
- Teacher profiles
- Subject allocation
- Attendance
- Timetable

### Examination Module
- Exams
- Marks entry
- Grading
- Report cards
- GPA calculation
- Rankings

### Finance Module
- Fee structures
- Invoices
- Payments
- Receipts
- Outstanding balances

### Communication
- SMS
- Email
- Notifications
- Announcements

---

## 5. Project Structure

```text
skuliApp/
│
├── frontend/            # React application
├── backend/             # Go API
├── docs/                # API docs, ERD, architecture
├── docker/
├── docker-compose.yml
├── .gitignore
└── README.md
```

### 5.1 Frontend (React)

```text
frontend/
│
├── public/
│
├── src/
│   ├── app/
│   │   ├── router/
│   │   ├── providers/
│   │   ├── store/
│   │   └── App.tsx
│   │
│   ├── assets/
│   │
│   ├── components/
│   │   ├── ui/
│   │   ├── tables/
│   │   ├── forms/
│   │   ├── charts/
│   │   ├── modals/
│   │   └── layout/
│   │
│   ├── features/
│   │   ├── auth/
│   │   ├── dashboard/
│   │   ├── students/
│   │   ├── teachers/
│   │   ├── parents/
│   │   ├── classes/
│   │   ├── subjects/
│   │   ├── attendance/
│   │   ├── timetable/
│   │   ├── exams/
│   │   ├── grading/
│   │   ├── finance/
│   │   ├── library/
│   │   ├── inventory/
│   │   ├── reports/
│   │   ├── notifications/
│   │   └── settings/
│   │
│   ├── hooks/
│   ├── services/
│   ├── lib/
│   ├── utils/
│   ├── types/
│   ├── constants/
│   └── main.tsx
│
├── package.json
└── vite.config.ts
```

### 5.2 Backend (Go)

```text
backend/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── config/
│
├── internal/
│   ├── auth/
│   ├── users/
│   ├── students/
│   ├── teachers/
│   ├── parents/
│   ├── classes/
│   ├── subjects/
│   ├── attendance/
│   ├── timetable/
│   ├── exams/
│   ├── grading/
│   ├── finance/
│   ├── library/
│   ├── inventory/
│   ├── reports/
│   ├── notifications/
│   └── settings/
│
├── pkg/
│   ├── database/
│   ├── middleware/
│   ├── logger/
│   ├── validator/
│   ├── jwt/
│   ├── response/
│   └── storage/
│
├── migrations/
├── docs/
├── uploads/
├── go.mod
└── go.sum
```

### 5.3 Module Structure

Every backend feature follows the same layout. Example — the `teachers` module:

```text
teachers/
│
├── handler.go        # HTTP request handling
├── service.go        # Business logic
├── repository.go     # Database access
├── model.go          # Database models
├── dto.go            # Request/response objects
├── routes.go         # Route registration
└── validator.go      # Input validation
```

### 5.4 Database Layer

```text
database/
├── migrations/
├── seeders/
└── factories/
```

---

## 6. Roles Reference

| Role         | Scope |
| ------------ | ----- |
| Super Admin  | Platform-wide administration |
| School Admin | Full control within a school |
| Principal    | School oversight |
| Teacher      | Classes, attendance, grading |
| Accountant   | Finance module |
| Parent       | View student information |
| Student      | Personal dashboard |

**Auth building blocks:** login · logout · refresh token · forgot password ·
change password · JWT middleware · role middleware.
