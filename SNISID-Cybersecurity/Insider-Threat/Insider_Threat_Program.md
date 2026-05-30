# SNISID Insider Threat Program

## 1. Objective
To detect and mitigate risks posed by authorized users (employees, contractors, partners) who may misuse their access to compromise national security.

## 2. Risk Categories
- **The Malicious Insider:** Intentionally steals data or sabotages systems for financial gain, ideology, or revenge.
- **The Careless Insider:** Accidentally causes a breach through negligence (e.g., clicking a phishing link, misconfiguring a bucket).
- **The Compromised Insider:** A legitimate user whose credentials have been stolen by an external actor.

## 3. Detection Matrix

| Risk | Detection Method | Indicators (IoC/IoB) |
| :--- | :--- | :--- |
| **Data Exfiltration** | DLP / SIEM | Mass uploads to cloud storage, large email attachments to personal accounts, USB usage. |
| **Privilege Abuse** | IAM Logs / SIEM | Accessing sensitive data not required for job role, frequent use of `sudo` or admin tools. |
| **Unusual Access** | Behavioral Analysis | Logins at 3 AM from a new location, access to systems never touched before. |
| **Mass Downloads** | File Server Logs | Downloading entire repositories or database dumps. |

## 4. Mitigation Strategies
- **Separation of Duties:** Ensuring no single person has full control over a critical process.
- **Dual Authorization:** Requiring two people to approve high-risk changes (e.g., PKI root key access).
- **Psychological Safety & Reporting:** Providing a safe way to report suspicious behavior (whistleblowing).
- **Mandatory Vacations:** Forcing employees to take leave to uncover ongoing fraudulent activities.

## 5. Privacy and Ethics
The program must balance security with privacy:
- **Transparency:** Users are notified that their activity on national systems is monitored.
- **Limited Access:** Only a small, vetted group of analysts can view insider threat alerts.
- **Legal Review:** All investigations must be approved by the Legal and HR departments.
