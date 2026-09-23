# TSC Template

**Reviewed**: 2026-09-23

Reference: [The Tech Stack Canvas](https://techstackcanvas.io/)

**Purpose**: the technology picture of an initiative — services, stack, integrations, infrastructure.
Generated as `docs/aics/<slug>/tsc.md` when the engineer invokes `/arcdlc:aic <slug> tsc` (flat
installs: `arcdlc-aic <slug> tsc`), alone or combined with another format (`arc42,tsc`).

**How to fill it**: keep the four groups and every heading below, in this order, with their colour
markers. Replace each description with the initiative's own content; a field that does not apply
stays in place with an explicit `Not applicable: <reason>`. Open the document with a
`# <Title>` heading and a one-line `> ` summary blockquote under it — `arctool sync` parses both for
the initiative registry. Prefer tables and short bullets over prose, and state only what the
interview, an ADR, or the code supports.

The **Ask** lines under a field are the questions `/arcdlc:aic` puts to the engineer to fill it, and
are not copied into the document. A field the interview deferred goes under Open questions, at the
end, with who answers it and by which phase.

## Business

### 🟠 Business Feature Description

Listing 3 main objectives that this project tries to reach, such as reducing costs by automating processes, reaching more customers, support a whole new business model. Example: Adjust salary per employee, Compare salaries, Prevent a pay gap, Less errors due to less manual steps. Link to AIC.

**Ask**: what does the problem cost today, and what does reaching these objectives save or earn?

### 🟠 System Name and Description

Describe a Service Purpose that is located inside Business Feature. Provide short overview.

### 🟠 Sizing Numbers

Listing some crucial numbers that give an idea of sizing and that influenced technology decisions. Expected number of users, requests per second, expected data volume. Example: 4-5 users at the same time, only at specific times of the year, very small data volume.

**Ask**: how many requests or user actions per second today and at the design target, and where do both
numbers come from? How many bytes per stored item, how many items alive at once, and for how long?

### 🟠 Major Quality Attributes

Quality attributes that influenced decisions. Examples are response times, adaptability, availability requirements. There are different approaches to document this. One is to list the ISO-25010 criterias that are most important. Example: Maintainability, Security, Reliability.

**Ask**: what measurement proves each attribute is met, and what number counts as a pass?

## Solution

### ⚪ Services

List of services inside the feature. Short overview of each service.

**Ask**: does each service have a reason that a package inside an existing service cannot meet?

### ⚪ Architecture Representation

Connection between services, DBs. C4 might be used.

### ⚪ Sequence Flow

A Sequence Diagram, if needed, that represents the technical flow between services and basic data.

**Ask**: how many calls to other services, and how many store reads and writes, does one user action make?
What happens when a step fails, times out, or arrives twice?

### ⚪ Frontend Technologies

The tools and frameworks used for developing the user interface and user experience, such as programming languages (e.g., JavaScript), libraries (e.g., React), and frameworks (e.g., Angular). If you are developing a mobile app, this is where you would list the technologies used, such as iOS, Android, SwiftUI, Flutter, etc.

### ⚪ Backend Technologies

The technologies used for server-side processing, data management, and business logic implementation, including programming languages (e.g., Python, Java), frameworks (e.g., Django, Ruby on Rails).

### ⚪ Data Storage and Management

The technologies used for data storage, retrieval, and processing, such as relational databases (e.g., MySQL), NoSQL databases (e.g., MongoDB), and data warehousing solutions (e.g., BigQuery).

**Ask**: what is the chosen store's throughput limit, where was it measured, and what happens when it is
reached? Which options lost, and by what criterion?

## Assessment

### 🟢 APIs and Integrations

The application programming interfaces (APIs) and third-party services used to extend the functionality of the product or service, such as payment gateways, email services, or social media integrations. This could also include integrations with internal solutions.

**Ask**: which partners can be slow, or absent, without warning? When one is, does the system fail open or
closed?

### 🟢 Security and Compliance

The tools, practices, and standards implemented to ensure the security and privacy of the application and its data, including encryption tools, security frameworks, and relevant regulations (e.g., GDPR, HIPAA).

**Ask**: who may call what, how is the caller authenticated, and which data is personal?

### 🟢 Testing and Quality Assurance

The tools and methodologies used to test the application's functionality, performance, and security, including automated testing frameworks (e.g., Selenium), continuous integration and continuous deployment (CI/CD) tools, and performance testing tools.

## Infrastructure

### ⚫ Infrastructure and Deployment

The platforms, tools, and services used for hosting, deploying, and managing the application, including cloud providers (e.g., AWS, Google Cloud), containerization tools (e.g., Docker), and deployment tools (e.g., Kubernetes).

**Ask**: what does it cost per month? How does it go live, and how fast does it go back?

### ⚫ Monitoring and Analytics

The tools and services used to monitor the application's performance, track user behavior, and gather insights for optimization and improvement, including application performance monitoring (APM) tools, log management solutions, and analytics platforms.

**Ask**: which metrics and alert values tell the person woken at 03:00 what broke?

### ⚫ Development Workflow and Collaboration

The tools and processes used to facilitate efficient development and collaboration among team members, such as version control systems (e.g., Git), project management tools (e.g., Jira), and communication platforms (e.g., Slack).

## Open questions

Every answer the interview deferred. A named unknown is fine; a hidden one is not.

| Question | Who answers | By when (phase or stage) |
| --- | --- | --- |
