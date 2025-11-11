# Improved MVP Backlog - Personal Finance Tracker

**Stack:** Svelte (Frontend) + Golang (Backend API) + PostgreSQL (Database)

---

## Epic 0 – Technical Foundation (NEW)

**Goal:** Establish core infrastructure before feature development.

### US 0.1 – Project Setup & Architecture

**As a developer**  
I want the basic project structure and tooling configured  
So that the team can develop features efficiently

**Estimate:** 5 points

**Technical Requirements:**
- Svelte app with SvelteKit for routing and SSR capabilities
- Golang API with proper project structure (handler/service/repository layers)
- PostgreSQL database with migration tool (golang-migrate or similar)
- Docker Compose for local development environment
- Environment configuration management (.env files)

**Acceptance Criteria:**

**Scenario 1 – Local development setup**
- Given a developer clones the repository
- When they run `docker-compose up` and follow setup instructions
- Then the Svelte app, Go API, and PostgreSQL are running and communicating

**Scenario 2 – API health check**
- Given the services are running
- When I call GET `/api/health`
- Then I receive a 200 response with service status

**Scenario 3 – Database migrations**
- Given migration files exist
- When I run the migration command
- Then the database schema is created/updated without errors

---

### US 0.2 – Authentication Infrastructure

**As a developer**  
I want JWT-based authentication infrastructure  
So that user sessions are secure and stateless

**Estimate:** 5 points

**Technical Requirements:**
- JWT token generation and validation in Golang
- Secure password hashing with bcrypt
- HTTP-only cookies for token storage
- Middleware for protected routes
- Token refresh mechanism

**Acceptance Criteria:**

**Scenario 1 – Token generation**
- Given valid user credentials
- When authentication succeeds
- Then a JWT token is generated with user claims and appropriate expiration

**Scenario 2 – Protected endpoints**
- Given an API endpoint requires authentication
- When a request is made without a valid token
- Then it returns 401 Unauthorized

**Scenario 3 – Token refresh**
- Given a user has a valid but expiring token
- When they request a token refresh
- Then a new token is issued without requiring re-login

---

## Epic 1 – Onboarding & Authentication (REVISED)

**Goal:** Let users securely sign up, sign in, and manage basic profile settings.

### US 1.1 – Email Sign-up

**As a new user**  
I want to create an account with email/password  
So that I can securely access my finance data

**Estimate:** 5 points (increased from 3)

**Technical Requirements:**
- PostgreSQL `users` table with indexes on email
- Golang endpoint: POST `/api/v1/auth/signup`
- Password validation: min 8 chars, uppercase, lowercase, number
- Email format validation with regex
- Bcrypt password hashing (cost factor 12)
- Rate limiting on signup endpoint (5 attempts per IP per hour)

**Database Schema:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    base_currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

**Acceptance Criteria:**

**Scenario 1 – Successful sign-up**
- Given I am on the sign-up page
- When I provide a valid email, strong password, and confirm password and click "Sign up"
- Then my account is created, I receive a JWT token, and I'm redirected to the dashboard

**Scenario 2 – Duplicate email**
- Given an account already exists with my email
- When I try to sign up with the same email
- Then I see an error "Email already registered" and the account is not created

**Scenario 3 – Validation errors**
- Given I am on the sign-up page
- When I enter an invalid email (e.g., "notanemail") or weak password (e.g., "123")
- Then I see specific validation messages and the account is not created

**Scenario 4 – Password confirmation mismatch**
- Given I am on the sign-up page
- When password and confirm password don't match
- Then I see "Passwords do not match" and cannot submit

---

### US 1.2 – Login & Logout (REVISED)

**As a registered user**  
I want to log in and log out  
So that I can securely access and end my session

**Estimate:** 4 points (increased from 3)

**Technical Requirements:**
- Golang endpoints: POST `/api/v1/auth/login`, POST `/api/v1/auth/logout`
- Set JWT in HTTP-only cookie
- Rate limiting on login (10 attempts per hour per IP)
- Clear cookie on logout
- Svelte store for authentication state

**Acceptance Criteria:**

**Scenario 1 – Successful login**
- Given I have a valid account
- When I enter my correct email and password and click "Log in"
- Then I receive a JWT token in an HTTP-only cookie and I'm taken to my dashboard

**Scenario 2 – Invalid credentials**
- Given I have a valid account
- When I enter an incorrect password
- Then I see an error "Invalid email or password" and remain on the login page

**Scenario 3 – Logout**
- Given I am logged in
- When I click "Log out"
- Then my JWT cookie is cleared and I'm redirected to the login page

**Scenario 4 – Session persistence**
- Given I am logged in
- When I close and reopen the browser
- Then I remain logged in if the token hasn't expired

---

### US 1.3 – Profile & Base Currency (REVISED)

**As a user**  
I want to set my display name and base currency  
So that amounts are shown consistently

**Estimate:** 3 points

**Technical Requirements:**
- Golang endpoint: PATCH `/api/v1/users/profile`
- Currency list: ISO 4217 codes (USD, EUR, GBP, JPY, etc.)
- Update `users` table
- Svelte form with currency dropdown

**Acceptance Criteria:**

**Scenario 1 – Update profile**
- Given I am logged in
- When I update my display name and save
- Then the new name appears in the header and profile page

**Scenario 2 – Set base currency**
- Given I am on my profile settings
- When I choose a base currency from a dropdown (e.g., EUR) and save
- Then all amounts are displayed with the correct currency symbol/code

**Scenario 3 – Invalid currency**
- Given I attempt to set an invalid currency code
- When I try to save
- Then I see an error and the currency is not changed

---

## Epic 2 – Transactions & Accounts (REVISED)

**Goal:** Let users record income/expenses and review them.

### US 2.1 – Manually Add Transaction (REVISED)

**As a user**  
I want to add expenses and income manually  
So that I can track my cash flow

**Estimate:** 5 points

**Technical Requirements:**
- PostgreSQL `transactions` table with indexes
- Golang endpoint: POST `/api/v1/transactions`
- Transaction types enum: 'income', 'expense'
- Foreign key to categories table
- Svelte form with date picker, category selector

**Database Schema:**
```sql
CREATE TYPE transaction_type AS ENUM ('income', 'expense');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type transaction_type NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    description VARCHAR(255),
    transaction_date DATE NOT NULL,
    category_id UUID REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
```

**Acceptance Criteria:**

**Scenario 1 – Add expense**
- Given I am logged in
- When I open "Add transaction", choose "Expense", enter amount (e.g., 50.00), date, category "Groceries", and description, then save
- Then the expense appears in my transaction list and total expenses increase by 50.00

**Scenario 2 – Add income**
- Given I am on "Add transaction"
- When I choose "Income", enter amount (e.g., 2000.00), date, category "Salary", and save
- Then the income appears in my transaction list and total income increases by 2000.00

**Scenario 3 – Validation**
- Given I leave amount empty or enter zero/negative
- When I try to save
- Then I see "Amount must be greater than 0" and the transaction is not saved

**Scenario 4 – Default values**
- Given I open "Add transaction"
- When the form loads
- Then today's date is pre-filled and my base currency is selected

---

### US 2.2 – View & Edit Transactions (REVISED)

**As a user**  
I want to see a list of my transactions and edit them  
So that I can correct mistakes

**Estimate:** 5 points (increased from 3)

**Technical Requirements:**
- Golang endpoints: 
  - GET `/api/v1/transactions?page=1&limit=50&sort=date_desc`
  - PATCH `/api/v1/transactions/{id}`
  - DELETE `/api/v1/transactions/{id}`
- Pagination with cursor or offset
- Svelte table/list component with sorting
- Optimistic UI updates for better UX

**Acceptance Criteria:**

**Scenario 1 – View list**
- Given I have created 100 transactions
- When I open the "Transactions" page
- Then I see the first 50 transactions sorted by date (newest first) with date, description, category, type, and amount

**Scenario 2 – Pagination**
- Given I have more than 50 transactions
- When I scroll to the bottom or click "Next"
- Then the next 50 transactions load

**Scenario 3 – Edit transaction**
- Given a transaction exists
- When I click "Edit", change the amount from 50 to 55, and save
- Then the transaction shows 55 and my totals are recalculated

**Scenario 4 – Delete transaction**
- Given a transaction exists
- When I click "Delete" and confirm in the modal
- Then the transaction is removed from the list and totals are updated

**Scenario 5 – Filter by type**
- Given I have both income and expense transactions
- When I select "Expenses only" filter
- Then only expense transactions are shown

---

### US 2.3 – Import Transactions from CSV (REVISED)

**As a user**  
I want to import a bank statement CSV  
So that I don't have to enter every transaction manually

**Estimate:** 13 points (increased from 8)

**Technical Requirements:**
- Golang endpoint: POST `/api/v1/transactions/import` (multipart/form-data)
- CSV parsing with encoding detection (UTF-8, ISO-8859-1)
- Column mapping UI in Svelte
- Duplicate detection: hash of (user_id, date, amount, description)
- Transaction creation in database transaction (all or nothing)
- Progress indicator for large files

**Expected CSV Format:**
```csv
Date,Description,Amount,Type
2024-01-15,Grocery Store,-45.50,expense
2024-01-16,Salary Deposit,2000.00,income
```

**Acceptance Criteria:**

**Scenario 1 – Successful import**
- Given I have a CSV with 50 valid transactions
- When I upload the file, map columns (date→Date, description→Description, amount→Amount, type→Type), and confirm
- Then all 50 transactions are created and appear in my list

**Scenario 2 – Column mapping**
- Given I upload a CSV with headers "fecha,concepto,importe"
- When I'm on the mapping screen
- Then I can map "fecha"→Date, "concepto"→Description, "importe"→Amount

**Scenario 3 – Invalid file**
- Given I upload a .txt file or CSV with missing required column
- When I attempt to import
- Then I see "Invalid file format. Expected CSV with Date, Amount columns" and no transactions are imported

**Scenario 4 – Duplicate handling**
- Given I import a CSV containing 10 transactions, 3 of which already exist (same date/amount/description)
- When I complete the import
- Then I see "7 transactions imported, 3 duplicates skipped" and only the 7 new ones are added

**Scenario 5 – Partial errors**
- Given a CSV with 10 rows, where row 5 has an invalid date format
- When I import
- Then I see "9 transactions imported, 1 failed: Row 5 - invalid date format" and the valid transactions are saved

**Scenario 6 – Large file handling**
- Given I upload a CSV with 5000 transactions
- When processing
- Then I see a progress bar and the import completes without timeout

---

## Epic 3 – Categories & Budgets (REVISED)

**Goal:** Help users structure spending and set budgets.

### US 3.1 – Manage Categories (REVISED)

**As a user**  
I want default categories and the ability to add my own  
So that I can group my spending meaningfully

**Estimate:** 4 points (increased from 3)

**Technical Requirements:**
- PostgreSQL `categories` table
- Seed data for default categories
- Golang endpoints: GET, POST, PATCH, DELETE `/api/v1/categories`
- Soft delete or reassignment for categories with transactions

**Database Schema:**
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7), -- hex color for UI
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- Default categories (user_id NULL or specific system user)
INSERT INTO categories (id, name, is_default) VALUES
    (gen_random_uuid(), 'Groceries', true),
    (gen_random_uuid(), 'Rent', true),
    (gen_random_uuid(), 'Transport', true),
    (gen_random_uuid(), 'Utilities', true),
    (gen_random_uuid(), 'Entertainment', true),
    (gen_random_uuid(), 'Healthcare', true),
    (gen_random_uuid(), 'Salary', true),
    (gen_random_uuid(), 'Other', true);
```

**Acceptance Criteria:**

**Scenario 1 – Default categories**
- Given I am a new user who just signed up
- When I open the categories page
- Then I see default categories: Groceries, Rent, Transport, Utilities, Entertainment, Healthcare, Salary, Other

**Scenario 2 – Add custom category**
- Given I am on the categories page
- When I click "Add Category", enter "Subscriptions", choose a color, and save
- Then "Subscriptions" appears in the list and is selectable when creating transactions

**Scenario 3 – Rename category**
- Given I have a custom category "Subscriptions"
- When I rename it to "Monthly Subscriptions" and save
- Then the category name is updated everywhere

**Scenario 4 – Delete category with transactions**
- Given a category "Subscriptions" has 10 transactions
- When I try to delete it
- Then I'm prompted to reassign those transactions to another category or confirm deletion with reassignment to "Other"

**Scenario 5 – Cannot delete default categories**
- Given "Groceries" is a default category
- When I try to delete it
- Then I see "Default categories cannot be deleted" or the delete button is disabled

---

### US 3.2 – Create Monthly Budgets per Category (REVISED)

**As a user**  
I want to set a monthly limit per category  
So that I can control my spending

**Estimate:** 5 points

**Technical Requirements:**
- PostgreSQL `budgets` table
- Golang endpoints: GET, POST, PATCH, DELETE `/api/v1/budgets`
- Unique constraint on (user_id, category_id, month, year)

**Database Schema:**
```sql
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    year INTEGER NOT NULL CHECK (year >= 2000),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, category_id, month, year)
);

CREATE INDEX idx_budgets_user_period ON budgets(user_id, year, month);
```

**Acceptance Criteria:**

**Scenario 1 – Create budget**
- Given I have categories
- When I set a budget of 500 for "Groceries" for the current month and save
- Then the budget appears in my budget list showing "Groceries: $500 for November 2024"

**Scenario 2 – Edit budget**
- Given a budget of 500 exists for "Groceries"
- When I change it to 600 and save
- Then budget calculations use 600

**Scenario 3 – No negative values**
- Given I try to create a budget
- When I enter -100 or 0
- Then I see "Budget must be greater than 0" and it's not saved

**Scenario 4 – Multiple budgets**
- Given I am setting budgets
- When I create budgets for 5 different categories for the same month
- Then all appear in my budget overview

**Scenario 5 – Duplicate prevention**
- Given I have a budget for "Groceries" in November 2024
- When I try to create another budget for "Groceries" in November 2024
- Then I see "Budget already exists for this category and month" and can edit the existing one

---

### US 3.3 – Track Budget vs Actual (REVISED)

**As a user**  
I want to see how much I've spent vs my budget  
So that I know if I'm on track

**Estimate:** 5 points

**Technical Requirements:**
- Golang endpoint: GET `/api/v1/budgets/summary?month=11&year=2024`
- SQL query aggregating transaction totals per category for the month
- Calculate: budgeted, spent, remaining, percentage used
- Svelte component with progress bars

**API Response Example:**
```json
{
  "month": 11,
  "year": 2024,
  "categories": [
    {
      "category_id": "uuid",
      "category_name": "Groceries",
      "budgeted": 500.00,
      "spent": 450.75,
      "remaining": 49.25,
      "percentage_used": 90.15,
      "is_overspent": false
    }
  ]
}
```

**Acceptance Criteria:**

**Scenario 1 – Budget overview**
- Given I have budgets for "Groceries" (500) and "Transport" (200) for November 2024
- And I've spent 450 on groceries and 180 on transport
- When I open the "Budgets" view for November
- Then I see:
  - Groceries: 450/500, 50 remaining, 90% used
  - Transport: 180/200, 20 remaining, 90% used

**Scenario 2 – Overspend highlighting**
- Given I have a budget of 500 for "Groceries"
- And I've spent 550
- When I view budgets
- Then "Groceries" shows in red or with a warning icon showing "-50 over budget"

**Scenario 3 – Filter by month**
- Given I have budgets for October, November, and December 2024
- When I select "October 2024"
- Then budget vs actual values for October are displayed

**Scenario 4 – No budget for category**
- Given I have transactions in "Entertainment" but no budget set
- When I view the budgets page
- Then "Entertainment" shows spent amount with "No budget set"

**Scenario 5 – Visual indicators**
- Given I am viewing budget progress
- When a category is under 50% used, it shows green
- When 50-80%, it shows yellow
- When 80-100%, it shows orange
- When over 100%, it shows red

---

## Epic 4 – Dashboard & Insights (REVISED)

**Goal:** Give users a clear overview of their financial situation and basic insights.

### US 4.1 – Overview Dashboard (REVISED)

**As a user**  
I want a dashboard with key metrics  
So that I can quickly understand my situation

**Estimate:** 8 points (increased from 5)

**Technical Requirements:**
- Golang endpoint: GET `/api/v1/dashboard?month=11&year=2024`
- Aggregate queries for totals
- Chart.js or D3.js for visualization (recommend Chart.js for simplicity)
- Svelte components: MetricCard, CategoryChart, RecentTransactions
- Server-side caching for dashboard data (Redis optional)

**Acceptance Criteria:**

**Scenario 1 – Current month summary**
- Given I have 3000 income and 2200 expenses in November 2024
- When I open the dashboard
- Then I see:
  - Total Income: $3,000
  - Total Expenses: $2,200
  - Net: +$800 (in green)

**Scenario 2 – Category chart**
- Given I have expenses: Groceries 500, Rent 1200, Transport 300, Entertainment 200
- When I view the dashboard
- Then I see a pie or donut chart showing:
  - Rent 54.5%
  - Groceries 22.7%
  - Transport 13.6%
  - Entertainment 9.1%

**Scenario 3 – Recent transactions widget**
- Given I have 50 transactions
- When I open the dashboard
- Then I see the 5 most recent transactions with date, description, category, and amount

**Scenario 4 – Empty state**
- Given I am a new user with no transactions
- When I open the dashboard
- Then I see a message "No transactions yet. Add your first transaction to get started" with a call-to-action button

**Scenario 5 – Month selector**
- Given I am viewing the dashboard
- When I change the month from November to October
- Then all metrics and charts update to show October's data

---

### US 4.2 – Monthly Trend Report (REVISED)

**As a user**  
I want to see spending trends over time  
So that I can understand how my behaviour changes

**Estimate:** 5 points (increased from 3)

**Technical Requirements:**
- Golang endpoint: GET `/api/v1/reports/trends?months=6`
- SQL query with GROUP BY year, month
- Line chart with Chart.js
- Svelte component with interactive tooltips

**Acceptance Criteria:**

**Scenario 1 – Trend chart**
- Given I have transactions from June to November 2024
- When I open "Reports" → "Monthly Trend"
- Then I see a line chart with two lines:
  - Total income per month
  - Total expenses per month
  - X-axis: months (Jun, Jul, Aug, Sep, Oct, Nov)
  - Y-axis: amount in base currency

**Scenario 2 – Month selection**
- Given I am viewing the trend chart
- When I click on "October" data point
- Then a modal or side panel shows October's breakdown by category

**Scenario 3 – Timeframe selection**
- Given I am viewing trends
- When I select "Last 12 months" instead of default 6
- Then the chart updates to show 12 months of data

**Scenario 4 – Net trend**
- Given the trend report
- When I toggle "Show Net Income"
- Then a third line appears showing (income - expenses) per month

---

### US 4.3 – Simple Overspending Insight (REVISED)

**As a user**  
I want to be warned when I'm close to or over budget  
So that I can adjust my behaviour

**Estimate:** 4 points (increased from 3)

**Technical Requirements:**
- Golang logic to check budget thresholds
- Warning threshold configurable in user settings (default 80%)
- Display warnings on dashboard and budget page
- Svelte notification component

**Acceptance Criteria:**

**Scenario 1 – Near budget warning**
- Given I have a budget of 500 for "Groceries"
- And my threshold is set to 80%
- When my spending reaches 400 (80%)
- Then I see a yellow warning indicator on the dashboard: "You've used 80% of your Groceries budget"

**Scenario 2 – Over budget alert**
- Given I have exceeded the budget for "Transport" (spent 250 of 200)
- When I view my dashboard or budgets page
- Then I see a red alert: "You've exceeded your Transport budget by $50"

**Scenario 3 – Multiple warnings**
- Given I have warnings for 3 categories
- When I view the dashboard
- Then I see a summary: "⚠️ 3 budget warnings" that expands to show all

**Scenario 4 – Configurable threshold**
- Given I am in settings
- When I change the warning threshold from 80% to 90%
- Then warnings only appear when I reach 90% of any budget

**Scenario 5 – Dismissible alerts**
- Given I see a budget warning on the dashboard
- When I click "Dismiss"
- Then it's hidden until the next day or until spending increases further

---

## Additional Recommendations

### Epic 5 – Data Export & Settings (NEW - Consider for post-MVP)

**US 5.1 – Export transactions to CSV**
- Users can download their data for backup or external analysis
- Estimate: 3 points

**US 5.2 – Account settings**
- Email change, password change, account deletion
- Estimate: 5 points

### Epic 6 – Multi-Currency Support (NEW - Consider for post-MVP)

**US 6.1 – Currency conversion**
- Store exchange rates
- Convert transactions to base currency for reporting
- Estimate: 8 points

---

## Technical Stack Implementation Notes

### Svelte Frontend
- **Routing:** SvelteKit with file-based routing
- **State Management:** Svelte stores for auth, user data
- **Forms:** Use `use:enhance` for progressive enhancement
- **Components:** Break down into reusable components (Button, Input, Card, Modal)
- **Charts:** Chart.js via svelte-chartjs wrapper
- **Date Picker:** svelte-flatpickr or native HTML5 date input

### Golang Backend
- **Framework:** Consider Chi, Gin, or Echo for routing
- **Structure:**
  ```
  /cmd/api         - main entry point
  /internal
    /handler       - HTTP handlers
    /service       - business logic
    /repository    - database queries
    /middleware    - auth, logging, cors
    /model         - data models
  ```
- **Validation:** go-playground/validator
- **Database:** pgx for PostgreSQL driver
- **Migrations:** golang-migrate or goose
- **Testing:** testify for assertions, httptest for API tests

### PostgreSQL
- **Connection Pooling:** Configure appropriate pool size (25-50 for MVP)
- **Indexing:** Ensure indexes on frequently queried columns
- **Backups:** Implement automated daily backups
- **Monitoring:** Log slow queries (>100ms)

### DevOps
- **Docker:** Separate containers for frontend, backend, database
- **CI/CD:** GitHub Actions for automated testing and deployment
- **Environment:** dev, staging, prod environments

---

## Updated Definition of Done

A user story is Done when:

### Functionality
- All acceptance criteria are met
- Feature works in latest Chrome, Firefox, Safari, Edge
- Responsive design (mobile: 375px+, tablet: 768px+, desktop: 1024px+)
- **API endpoint documented in OpenAPI/Swagger spec**
- **Database migrations applied and tested**

### Quality
- **Golang unit tests:** >70% coverage for service layer
- **Integration tests:** Key API flows tested with httptest
- **Svelte component tests:** Critical components have tests
- No known critical or high-severity bugs
- **Manual testing checklist completed**

### Security & Compliance
- Input validation on both frontend and backend
- SQL injection prevention (use parameterized queries)
- XSS prevention (Svelte auto-escapes by default)
- HTTPS enforced in production
- Passwords hashed with bcrypt (cost 12)
- JWT tokens with appropriate expiration (15 min access, 7 day refresh)
- Rate limiting implemented on auth and import endpoints
- CORS configured correctly
- No sensitive data in logs or error messages

### UX & Content
- UI matches approved designs/wireframes
- Loading states for async operations
- Error messages are clear and actionable
- Form validation shows inline errors
- Success feedback (toasts, messages)
- Keyboard navigation works for forms
- WCAG 2.1 AA contrast ratios
- Semantic HTML with ARIA labels where needed

### Performance & Ops
- API responses <200ms for queries, <1s for complex reports
- Database queries optimized with EXPLAIN ANALYZE
- No N+1 query problems
- Frontend bundle size <500KB
- Lazy loading for heavy components
- Logging for auth events, errors, slow queries
- Error tracking (consider Sentry)
- Feature deployed to target environment
- Release notes updated

### Documentation
- API endpoints documented with examples
- Database schema changes noted in migration
- Complex business logic has comments
- User-facing help text or tooltips where needed
- README updated if setup changes

---

## Recommended Development Order

1. **Epic 0:** Technical Foundation (US 0.1, 0.2)
2. **Epic 1:** Authentication (US 1.1, 1.2, 1.3)
3. **Epic 2:** Core Transactions (US 2.1, 2.2)
4. **Epic 3:** Categories (US 3.1)
5. **Epic 4:** Dashboard (US 4.1 - basic version)
6. **Epic 3:** Budgets (US 3.2, 3.3)
7. **Epic 4:** Insights (US 4.2, 4.3)
8. **Epic 2:** CSV Import (US 2.3) - complex, do after core features work

---

## Estimated Timeline

**Total Points:** ~93 points

If 1 developer, velocity ~10-15 points/sprint (2 weeks):
- **Optimistic:** 6-7 sprints (12-14 weeks / 3-3.5 months)
- **Realistic:** 8-10 sprints (16-20 weeks / 4-5 months)

This includes time for testing, bug fixes, and learning the stack.

---

## Risk Mitigation

1. **CSV Import Complexity (US 2.3)**
   - Risk: High - many edge cases
   - Mitigation: Build a robust parser, extensive testing with real bank CSVs

2. **Authentication Security**
   - Risk: Medium - security vulnerabilities
   - Mitigation: Use proven libraries (bcrypt, JWT), security audit before launch

3. **Database Performance**
   - Risk: Medium - slow queries as data grows
   - Mitigation: Proper indexing, query optimization, pagination

4. **Scope Creep**
   - Risk: High - feature requests during development
   - Mitigation: Stick to MVP, maintain backlog for post-MVP features

5. **Third-party Dependencies**
   - Risk: Low-Medium - library updates/deprecations
   - Mitigation: Pin versions, monitor security advisories

---

## Success Metrics (Define and Track)

### User Adoption Metrics
- **New user signups per week:** Target 50+ in first month
- **Active users (7-day):** Users who log in and view transactions
- **User retention:** % of users who return after 7 days, 30 days
- **Time to first transaction:** How long after signup do users add their first transaction

### Feature Usage Metrics
- **Transactions created:** Average per user per week (target: 10+)
- **CSV imports:** % of users who use import feature
- **Budget creation:** % of users who set at least one budget
- **Dashboard visits:** Daily active users viewing dashboard
- **Category customization:** % of users who create custom categories

### Technical Metrics
- **API response time:** P95 <200ms for queries, <1s for reports
- **Error rate:** <1% of API requests fail
- **Uptime:** 99.5% availability
- **Page load time:** <2s for initial dashboard load

### Business/Value Metrics
- **User satisfaction:** NPS score or simple feedback survey
- **Feature requests:** Track top 5 requested features
- **Bug reports:** Critical bugs resolved within 24h

---

## Testing Strategy

### Unit Tests (Golang Backend)
```go
// Example test structure
func TestCreateTransaction(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateTransactionRequest
        want    error
        wantErr bool
    }{
        {
            name: "valid expense transaction",
            input: CreateTransactionRequest{
                Type: "expense",
                Amount: 50.00,
                Date: "2024-11-11",
                CategoryID: "uuid",
            },
            wantErr: false,
        },
        {
            name: "invalid negative amount",
            input: CreateTransactionRequest{
                Amount: -50.00,
            },
            wantErr: true,
        },
    }
    // Run tests...
}
```

**Coverage targets:**
- Service layer: >80%
- Repository layer: >70%
- Handlers: >60%

### Integration Tests (Golang)
```go
func TestTransactionFlow(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Test full flow
    // 1. Create user
    // 2. Login
    // 3. Create transaction
    // 4. Verify in database
    // 5. Retrieve via API
}
```

### Component Tests (Svelte)
```javascript
// Example with @testing-library/svelte
import { render, fireEvent } from '@testing-library/svelte';
import TransactionForm from './TransactionForm.svelte';

test('shows validation error for empty amount', async () => {
    const { getByText, getByLabelText } = render(TransactionForm);
    
    const submitButton = getByText('Save');
    await fireEvent.click(submitButton);
    
    expect(getByText('Amount is required')).toBeInTheDocument();
});
```

### E2E Tests (Playwright)
```javascript
test('user can create and view transaction', async ({ page }) => {
    // Login
    await page.goto('/login');
    await page.fill('[name="email"]', 'test@example.com');
    await page.fill('[name="password"]', 'Password123!');
    await page.click('button[type="submit"]');
    
    // Create transaction
    await page.click('text=Add Transaction');
    await page.fill('[name="amount"]', '50.00');
    await page.selectOption('[name="category"]', 'Groceries');
    await page.click('text=Save');
    
    // Verify appears in list
    await expect(page.locator('text=50.00')).toBeVisible();
});
```

**E2E Test Coverage:**
- Critical user flows (signup, login, create transaction, view dashboard)
- Budget creation and tracking
- CSV import happy path
- Run on CI before deployment

---

## Security Checklist

### Authentication & Authorization
- [ ] Passwords hashed with bcrypt (cost ≥12)
- [ ] JWT tokens with short expiration (15 min access token)
- [ ] Refresh token mechanism implemented
- [ ] HTTP-only cookies for token storage
- [ ] CSRF protection for state-changing operations
- [ ] Rate limiting on auth endpoints (5 signup, 10 login per hour per IP)
- [ ] Account lockout after N failed login attempts
- [ ] Secure password reset flow (if implemented)

### Input Validation
- [ ] All user inputs validated on backend (never trust frontend)
- [ ] Email format validation
- [ ] Amount validation (positive numbers only)
- [ ] SQL injection prevention (parameterized queries only)
- [ ] File upload validation (CSV only, max 10MB)
- [ ] XSS prevention (Svelte auto-escapes, but verify user-generated content)

### API Security
- [ ] CORS configured (whitelist frontend domain only)
- [ ] HTTPS enforced in production
- [ ] Security headers (X-Content-Type-Options, X-Frame-Options, etc.)
- [ ] Request size limits
- [ ] Rate limiting on expensive endpoints (imports, reports)
- [ ] API versioning (/api/v1/)

### Data Protection
- [ ] User data isolated (all queries filter by user_id)
- [ ] Foreign key constraints prevent orphaned data
- [ ] Sensitive data not logged (passwords, tokens)
- [ ] Database connection string not in code (environment variables)
- [ ] Regular automated backups
- [ ] Data retention policy defined

### Third-party Dependencies
- [ ] Regular dependency updates
- [ ] Security audit with `npm audit` / `go list -m all | nancy sleuth`
- [ ] Pin versions in package.json / go.mod
- [ ] Review licenses for compliance

---

## Performance Optimization Guidelines

### Database Optimization
```sql
-- Example: Optimize transaction listing query
CREATE INDEX idx_transactions_user_date 
ON transactions(user_id, transaction_date DESC);

-- Analyze query performance
EXPLAIN ANALYZE
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY transaction_date DESC
LIMIT 50;

-- Consider materialized view for dashboard
CREATE MATERIALIZED VIEW user_monthly_summary AS
SELECT 
    user_id,
    DATE_TRUNC('month', transaction_date) as month,
    type,
    SUM(amount) as total
FROM transactions
GROUP BY user_id, month, type;

-- Refresh periodically or on transaction insert
```

### Backend Optimization
- Use connection pooling (pgx pool)
- Implement caching for dashboard data (Redis optional)
- Use prepared statements for repeated queries
- Batch inserts for CSV import (1000 rows at a time)
- Pagination for all list endpoints
- Database queries run in parallel where possible (goroutines)

### Frontend Optimization
```javascript
// Lazy load heavy components
<script>
    import { onMount } from 'svelte';
    let ChartComponent;
    
    onMount(async () => {
        const module = await import('./HeavyChart.svelte');
        ChartComponent = module.default;
    });
</script>

// Virtual scrolling for long transaction lists
import { VirtualList } from 'svelte-virtual-list';

// Debounce search inputs
import { debounce } from 'lodash-es';
const handleSearch = debounce((query) => {
    // API call
}, 300);
```

### Monitoring
- Log slow queries (>100ms)
- Monitor API endpoint response times
- Track database connection pool usage
- Set up alerts for error rate spikes
- Monitor disk space usage

---

## Deployment Architecture

### Recommended Setup

```
[Users] → [Cloudflare/CDN] → [Load Balancer]
                                    ↓
                          [Svelte App (Static)]
                                    ↓
                          [Golang API Server(s)]
                                    ↓
                          [PostgreSQL Database]
                                    ↓
                          [Backup Storage]
```

### Environment Configuration

**Development:**
- Docker Compose with hot reload
- Local PostgreSQL
- Mock external services

**Staging:**
- Mirrors production setup
- Separate database
- Used for final testing before release

**Production:**
- Managed PostgreSQL (AWS RDS, DigitalOcean, etc.)
- Multiple API server instances (if needed)
- Automated backups (daily minimum)
- SSL/TLS certificates (Let's Encrypt)
- Environment variables for secrets
- Health check endpoints

### CI/CD Pipeline (GitHub Actions Example)

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
      - name: Check coverage
        run: |
          go tool cover -func=coverage.out
          
  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      - name: Run tests
        run: npm test
      - name: Build
        run: npm run build
        
  deploy-staging:
    needs: [test-backend, test-frontend]
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to staging
        run: |
          # Your deployment commands
          
  deploy-production:
    needs: [test-backend, test-frontend]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to production
        run: |
          # Your deployment commands
```

---

## API Documentation Examples

### Authentication Endpoints

#### POST /api/v1/auth/signup
**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "display_name": null,
    "base_currency": "USD"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "password",
      "message": "Password must be at least 8 characters"
    }
  ]
}
```

#### POST /api/v1/auth/login
**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "display_name": "John Doe",
    "base_currency": "USD"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Transaction Endpoints

#### GET /api/v1/transactions
**Query Parameters:**
- `page` (int, default: 1)
- `limit` (int, default: 50, max: 100)
- `type` (string, optional: "income" | "expense")
- `category_id` (UUID, optional)
- `start_date` (date, optional: YYYY-MM-DD)
- `end_date` (date, optional: YYYY-MM-DD)
- `sort` (string, default: "date_desc", options: "date_asc" | "date_desc" | "amount_asc" | "amount_desc")

**Response (200 OK):**
```json
{
  "transactions": [
    {
      "id": "650e8400-e29b-41d4-a716-446655440000",
      "type": "expense",
      "amount": 45.50,
      "currency": "USD",
      "description": "Grocery shopping",
      "transaction_date": "2024-11-11",
      "category": {
        "id": "750e8400-e29b-41d4-a716-446655440000",
        "name": "Groceries",
        "color": "#4CAF50"
      },
      "created_at": "2024-11-11T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total_items": 247,
    "total_pages": 5
  }
}
```

#### POST /api/v1/transactions
**Request:**
```json
{
  "type": "expense",
  "amount": 45.50,
  "currency": "USD",
  "description": "Grocery shopping",
  "transaction_date": "2024-11-11",
  "category_id": "750e8400-e29b-41d4-a716-446655440000"
}
```

**Response (201 Created):**
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "type": "expense",
  "amount": 45.50,
  "currency": "USD",
  "description": "Grocery shopping",
  "transaction_date": "2024-11-11",
  "category": {
    "id": "750e8400-e29b-41d4-a716-446655440000",
    "name": "Groceries",
    "color": "#4CAF50"
  },
  "created_at": "2024-11-11T10:30:00Z"
}
```

### Dashboard Endpoint

#### GET /api/v1/dashboard
**Query Parameters:**
- `month` (int, 1-12, default: current month)
- `year` (int, default: current year)

**Response (200 OK):**
```json
{
  "period": {
    "month": 11,
    "year": 2024
  },
  "summary": {
    "total_income": 3000.00,
    "total_expenses": 2245.75,
    "net": 754.25,
    "currency": "USD"
  },
  "spending_by_category": [
    {
      "category_id": "750e8400-e29b-41d4-a716-446655440000",
      "category_name": "Groceries",
      "total": 450.75,
      "percentage": 20.07,
      "color": "#4CAF50"
    },
    {
      "category_id": "850e8400-e29b-41d4-a716-446655440000",
      "category_name": "Rent",
      "total": 1200.00,
      "percentage": 53.43,
      "color": "#2196F3"
    }
  ],
  "recent_transactions": [
    {
      "id": "650e8400-e29b-41d4-a716-446655440000",
      "type": "expense",
      "amount": 45.50,
      "description": "Grocery shopping",
      "transaction_date": "2024-11-11",
      "category_name": "Groceries"
    }
  ],
  "budget_alerts": [
    {
      "category_id": "750e8400-e29b-41d4-a716-446655440000",
      "category_name": "Groceries",
      "budgeted": 500.00,
      "spent": 450.75,
      "percentage_used": 90.15,
      "alert_type": "warning"
    }
  ]
}
```

---

## Database Migration Examples

### Migration 001 - Initial Schema

**001_create_users_table.up.sql:**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    base_currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

**001_create_users_table.down.sql:**
```sql
DROP TABLE IF EXISTS users;
```

### Migration 002 - Categories and Transactions

**002_create_categories_transactions.up.sql:**
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, name)
);

CREATE TYPE transaction_type AS ENUM ('income', 'expense');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type transaction_type NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    description VARCHAR(255),
    transaction_date DATE NOT NULL,
    category_id UUID REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_user_type ON transactions(user_id, type);

-- Insert default categories
INSERT INTO categories (name, is_default, color) VALUES
    ('Groceries', true, '#4CAF50'),
    ('Rent', true, '#2196F3'),
    ('Transport', true, '#FF9800'),
    ('Utilities', true, '#9C27B0'),
    ('Entertainment', true, '#E91E63'),
    ('Healthcare', true, '#00BCD4'),
    ('Salary', true, '#8BC34A'),
    ('Other', true, '#607D8B');
```

---

## Error Handling Standards

### Backend Error Responses

```go
// Standardized error response structure
type ErrorResponse struct {
    Error   string            `json:"error"`
    Details []ValidationError `json:"details,omitempty"`
    Code    string            `json:"code,omitempty"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// Error codes
const (
    ErrCodeValidation      = "VALIDATION_ERROR"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeDuplicate       = "DUPLICATE"
    ErrCodeServerError     = "INTERNAL_SERVER_ERROR"
    ErrCodeRateLimit       = "RATE_LIMIT_EXCEEDED"
)
```

### Frontend Error Handling

```javascript
// Centralized API error handler
export async function apiRequest(url, options = {}) {
    try {
        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
        });

        if (!response.ok) {
            const error = await response.json();
            throw new ApiError(error.error, error.code, error.details);
        }

        return await response.json();
    } catch (error) {
        if (error instanceof ApiError) {
            throw error;
        }
        // Network or other errors
        throw new ApiError('Network error. Please try again.', 'NETWORK_ERROR');
    }
}

// Custom error class
class ApiError extends Error {
    constructor(message, code, details = []) {
        super(message);
        this.code = code;
        this.details = details;
    }
}

// Usage in Svelte component
async function handleSubmit() {
    try {
        const result = await apiRequest('/api/v1/transactions', {
            method: 'POST',
            body: JSON.stringify(formData),
        });
        showSuccess('Transaction created successfully');
    } catch (error) {
        if (error.code === 'VALIDATION_ERROR') {
            // Show field-specific errors
            validationErrors = error.details;
        } else {
            // Show general error toast
            showError(error.message);
        }
    }
}
```

---

## Go Ahead and Build! 🚀

This improved backlog provides:

1. ✅ **Technical specificity** for your Svelte + Golang + PostgreSQL stack
2. ✅ **Database schemas** with proper indexing
3. ✅ **API endpoint specifications** with examples
4. ✅ **Security considerations** at every level
5. ✅ **Testing strategies** for all layers
6. ✅ **Performance optimization** guidelines
7. ✅ **Deployment architecture** recommendations
8. ✅ **Realistic estimates** with increased complexity where needed
9. ✅ **Epic 0** for technical foundation (often overlooked!)
10. ✅ **Additional recommendations** for post-MVP features

**Key Improvements Over Original:**
- Added Epic 0 for infrastructure setup
- Increased estimates for complex features (CSV import, dashboard)
- Added technical requirements to every user story
- Provided database schemas with proper constraints
- Included API documentation examples
- Enhanced Definition of Done with stack-specific items
- Added security checklist and performance guidelines

**Next Steps:**
1. Review and adjust estimates based on your experience level
2. Set up the development environment (Epic 0)
3. Start with Epic 1 (Authentication) - foundation for everything
4. Build iteratively, deploying to staging after each epic
5. Gather user feedback early and often

Good luck with your finance tracker MVP! 💰📊
