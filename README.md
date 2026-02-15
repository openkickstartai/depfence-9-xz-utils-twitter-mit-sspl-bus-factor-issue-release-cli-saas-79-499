# DepFence 🛡️

**Dependency health scoring & supply chain risk monitoring engine.**

Stop learning about abandoned deps from Twitter. DepFence scores every dependency on maintainer activity, bus factor, and license risk.

Go single binary. Zero config. CI-ready.

## 🚀 Quick Start

```bash
# Install
go install github.com/depfence/depfence@latest

# Auto-detect and scan
depfence

# Full scoring with GitHub API
export GITHUB_TOKEN=ghp_xxx
depfence -f go.mod

# CI gate: fail if any dep scores below 40
depfence -min-score 40

# Machine-readable output
depfence -format json
depfence -format csv > report.csv
```

## 📊 What It Scores

| Signal | Measurement | Risk |
|--------|------------|------|
| Maintainer Activity | Days since last push | Abandoned projects |
| Bus Factor | Active contributor count | Single-maintainer risk |
| License Risk | SPDX classification | SSPL/BSL surprises |

## 💰 Pricing

| Feature | Free (CLI) | Pro $79/mo | Enterprise $499/mo |
|---------|-----------|-----------|--------------------|
| Go/npm/pip scanning | ✅ | ✅ | ✅ |
| Health scoring (rate-limited) | ✅ | Unlimited | Unlimited |
| JSON/CSV export | ✅ | ✅ | ✅ |
| CI/CD exit codes | ✅ | ✅ | ✅ |
| Transitive dep analysis | ❌ | ✅ | ✅ |
| License change alerts | ❌ | ✅ | ✅ |
| Slack/PagerDuty integration | ❌ | ✅ | ✅ |
| Historical trend dashboard | ❌ | ✅ | ✅ |
| SOC2/ISO27001 PDF reports | ❌ | ❌ | ✅ |
| SSO/SAML | ❌ | ❌ | ✅ |
| Multi-repo monitoring | ❌ | 10 repos | Unlimited |
| Custom policy engine | ❌ | ❌ | ✅ |

## 🤔 Why Pay?

- **One abandoned dep costs $50K+** in emergency migration. DepFence catches it 6 months early.
- **License violations** block M&A, IPO, SOC2 audits. Continuous monitoring = audit evidence on demand.
- **xz-utils proved it**: supply chain attacks are real. $79/mo is insurance against a $500K incident.
- **ROI**: $79/mo vs one incident response = **100x+ return**.

## Supported Ecosystems

✅ Go (go.mod) · ✅ Node.js (package.json) · ✅ Python (requirements.txt) · 🔜 Rust · 🔜 Java

## License

BSL 1.1 — Free for evaluation. Commercial license required for production.
