# Scenario: Incident Triage

1. Planner agent gathers context and decomposes incident scope.
2. RAG tool fetches similar incidents by vector similarity.
3. Coder agent drafts patch and rollback plan.
4. Reviewer agent checks blast radius and policy compliance.
5. HITL gate requests human approval if confidence < threshold.
6. Final summary is piped to deployment owner.
