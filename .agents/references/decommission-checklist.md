# Decommission Checklist

1. Confirm the target repository or runtime surface and the end state.
2. Inventory the live surfaces that retirement must remove, including any Swarm stack, workflow, or deploy contract that would otherwise keep the target runnable.
3. Update the management ledgers that track active, sleepy, or archived repositories.
4. Record the retirement in the durable event log if the repo is managed there.
5. Remove or retire the live deployment surface before marking the retirement complete.
6. Preserve only the context that still teaches or supports future audits.
7. Close any related campaign or issue records.
8. Verify that the active pipeline no longer treats the repository as live work and that the retired runtime surface is absent from the cluster or deployment target.
