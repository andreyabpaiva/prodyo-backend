# Prodyo Backend API
## Features

- **Project Management**: Create, read, update, and delete projects with team members
- **Iteration Tracking**: Manage development iterations with tasks and metrics
- **Productivity Metrics**: Track speed, rework, and instability indicators
- **Quality Tracking**: Monitor bugs and improvements per task
- **Action Planning**: Create causes and actions based on productivity analysis
- **Iteration analysis**: Generate iteration analysis based on indicators ranges and iteration performance
- **Automatic Migrations**: Database migrations run automatically on startup
- **Docker Support**: Easy deployment with Docker Compose

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Running with Docker

```bash
docker-compose up
```

The API will be available at `http://localhost:8080`

## Class Diagram

```mermaid
classDiagram
    class User {
        - id: UUID
        - name: string
        - email: string
        - passwordHash: string
        - projectID: UUID
        + getById(id UUID): Project
        + update(user User): void
        + create(user User): void
    }

    class Project {
        - id: UUID
        - name: string
        - description: string
        - members: User[]
        - color: string
        - prodRange: ProductivityRange
        - createdAt: Time
        - updatedAt: Time
        + getAll(userId UUID): Project[]
        + getById(id UUID): Project
        + delete(id UUID): void
        + update(project Project): void
        + create(project Project): void
    }

    class ProductivityRange {
        - ok: int[]
        - alert: int[]
        - critical: int[]
    }

    class Iteration {
        - id: UUID
        - projectId: UUID
        - number: int
        - description: string
        - startAt: Time
        - endAt: Time
        - tasks: Task[]
        + getAll(projId UUID): Iteration[]
        + getById(id UUID): Iteration
        + delete(id UUID): void
        + create(iteration Iteration): void
    }

    class Indicator {
        - id: UUID
        - iterationId: UUID
        - causes: Cause[]
        - actions: Action[]
        + get(iterationId UUID): Indicator
    }

    class Task {
        - id: UUID
        - iterationId: UUID
        - name: string
        - description: string
        - assignee: User
        - status: StatusEnum
        - timer: Time
        - tasks: Task[]
        - Improvements: Improv[]
        - Bugs: Bug[]
        + getAll(iterationId UUID): Task[]
        + getById(id UUID): Task
        + delete(id UUID): void
        + create(task Task): void
        + update(task Task): void
    }

    class Improv {
        - id: UUID
        - taskId: UUID
        - assignee: User
        - number: int
        - description: string
        + getAll(taskId UUID): Improv[]
        + getById(id UUID): Improv
        + create(improv Improv): void
    }

    class Bug {
        - id: UUID
        - taskId: UUID
        - assignee: User
        - number: int
        - description: string
        + getAll(taskId UUID): Bug[]
        + getById(id UUID): Bug
        + create(bug Bug): void
    }

    class Cause {
        - id: UUID
        - indicatorId: UUID
        - metric: MetricEnum
        - description: string
        - productivityLevel: ProductivityEnum
        + get(indicatorId UUID): Cause
    }

    class Action {
        - id: UUID
        - indicatorId: UUID
        - description: string
        - cause: Cause
        + get(indicatorId UUID): Action
    }

    class ProductivityEnum {
        <<Enum>>
        Ok
        Alert
        Critical
    }

    class MetricEnum {
        <<Enum>>
        WorkVelocity
        ReworkIndex
        InstabilityIndex
    }

    class StatusEnum {
        <<Enum>>
        NotStarted
        InProgress
        Completed
    }

    Project "1" --> "1" ProductivityRange
    Project "1" --> "0..*" User
    Project "1" --> "1..*" Iteration
    Iteration "1" --> "0..*" Task
    Iteration "1" --> "1" Indicator
    Task "1" --> "0..*" Improv
    Task "1" --> "0..*" Bug
    Indicator "1" --> "0..*" Cause
    Indicator "1" --> "0..*" Action
    Action "1" --> "1" Cause
    Cause --> ProductivityEnum
    Cause --> MetricEnum
    Task --> StatusEnum
```

## License

MIT
