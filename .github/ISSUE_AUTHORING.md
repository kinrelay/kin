# Issue Authoring Policy

本文件定義 Kin GitHub issue 的 authoring contract。Issue 是目前 task 的最高優先 source of truth，但它應只承載 **task-specific authority**；可重用的 repository policy 與 procedure 必須留在 canonical `AGENTS.md`、local `AGENTS.md` 與 `.agents/skills/`。

## Issue 應包含什麼

Issue 應清楚描述這個 task 自己的：

- Goal / problem statement
- task-specific context
- Acceptance Criteria
- Explicit Non-goals
- dependencies / links
- 只有此 task 才成立的 assumptions、constraints 或 exceptions

這些內容決定 agent 在該 issue 可做什麼、完成條件是什麼；本 policy **不降低 current issue 在 authorized scope 與 AC 上的最高優先權**。

## Issue 不應複製什麼

不要把下列 reusable policy 整段 copy 進 issue：

- global architecture / dependency-direction rules
- generic testing / TDD / review / merge workflow
- generic coding-agent authority instructions
- reusable privacy / provider policy
- generic PR evidence / review-disposition procedure
- 已由 root/local `AGENTS.md` 或 repo-local skill 定義的其他 repository-wide contract

原因不是節省篇幅而已。因為 current issue precedence 高於 repository contracts，過時的 duplicated policy 可能意外覆蓋較新的 canonical rule，造成 instruction drift。

若 task 確實需要例外，issue 應只寫出 **task-specific exception**，並明確說明它與 canonical rule 的差異；不要為了表達例外而複製整份 policy。

## 建議 Agent reference pattern

```md
## Agent references

Follow the current repository contracts and the skills relevant to this task:

- root `AGENTS.md`
- nearest applicable local `AGENTS.md`
- `.github/AGENT_COORDINATION.md`
- relevant `.agents/skills/*/SKILL.md`

Task-specific constraints:
- <only constraints unique to this issue>
```

只有真正與 task 有關的 reference 才需要列出；不要把 referenced 文件內容再次貼進 issue。

## Existing issues

不要求批次重寫所有歷史 issue。當既有 issue 被重新使用、修改，或 duplicated policy 已形成明顯 drift 風險時，再在該 task scope 內精簡即可。
