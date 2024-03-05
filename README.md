# Домашка к курсу ["Асинхронная архитектура" by ToughDevSchool](https://tough-dev.school/architecture)


## Event Storming

### Events By Context

#### Task Tracker

| Action                 | Actor            | Command          | Data                                    | Event           | Description |
|------------------------|------------------|------------------|-----------------------------------------|-----------------|-------------|
| Create a Task          | Anyone           | CreateTask       | `Task{Description, Status, Assignee}`   | Task.Created    |             |
| Resolve Task           | Anyone           | ChangeTaskStatus | `Task{TaskID, Status}`                  | Task.Resolved   |             |
| Assign Task            | Admin or Manager | AssignTasks      | `Task{TaskID, Status: Open, Assignee}`  | []Task.Assigned |             |
| (RM) Send Notification | Task.Changed     | -                | -                                       | -               |             |
    
#### Accounting

| Action                       | Actor                | Command       | Data                         | Event                   | Description                                            |
|------------------------------|----------------------|---------------|------------------------------|-------------------------|--------------------------------------------------------|
| Withdraw Money               | Task.Assigned        | WithdrawMoney | `Account{UserID, MoneyDiff}` | Account.BalanceChanged  |                                                        |
| Deposit Money                | Task.Resolved        | DepositMoney  | `Account{UserID, MoneyDiff}` | Account.BalanceChanged  |                                                        |
| Reset Balance                | Cron                 | ResetBalance  |                              | Account.BalanceReset    | Reset balance all of the workers at the end of the day |
| Commit Payout                | Account.BalanceReset | CommitPayout  | `Account{UserID}`            | Account.PayoutCommitted |                                                        |
| (RM) Get Worker's Dashboard  | Worker               | -             | -                            | -                       |                                                        |
| (RM) Get Daily Balance       | Manager              |               |                              |                         |                                                        |
| (RM) Send Daily Notification | Cron                 |               |                              |                         |                                                        |

#### Analytics

| Action                                                  | Actor       | Command | Data | Event | Description                                                                                                                                          |
|---------------------------------------------------------|-------------|---------|------|-------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| (RM) Get Daily Balance                                  | Top-Manager |         |      |       |                                                                                                                                                      |
| (RM) Get The Most Expensive Task For [Day, Month, Year] | Top-Manager |         |      |       | 03.03 - самая дорогая задача - 28$; 02.03 - самая дорогая задача - 38$; 01.03 - самая дорогая задача - 23$; 01-03 марта - самая дорогая задача - 38$ |

## Domains

![Domains](./docs/domains.png)