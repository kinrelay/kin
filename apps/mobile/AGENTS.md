# Mobile Agent Contract

This contract applies to work under `apps/mobile/` and inherits the repository root `AGENTS.md`.

## Scope

`apps/mobile` is the Expo / React Native / TypeScript application. It follows the same interaction-first and boundary-first philosophy as the backend without forcing unnecessary DDD ceremony into presentation-only code.

## Required boundaries

- Start from the user interaction, screen behavior, and state transitions before choosing component structure or integration details.
- Keep presentation/UI state, feature orchestration, meaningful domain/application policy, and adapters conceptually distinct.
- API clients, native SDKs, device capabilities, storage, analytics, and external providers are outer boundaries.
- Components must not absorb business policy merely because the UI is the easiest place to implement it.
- Shared business rules should live in an appropriate domain/application layer or shared contract instead of being duplicated across screens.
- Presentation-only behavior does not require artificial aggregates, entities, or value objects when there is no meaningful domain invariant.

## Design order

For mobile features, prefer:

1. User interaction and expected behavior.
2. State transitions and policy boundaries.
3. Tests for meaningful behavior.
4. Application/API/native ports.
5. Components, navigation, networking, storage, and native/provider adapters.

Do not begin feature design from component trees, API calls, native SDK usage, or state-management library mechanics.

## Extensibility

Keep this file minimal. Expo/React Native conventions, navigation patterns, state-management guidance, native capability procedures, EAS workflows, accessibility rules, and mobile testing commands should move to a dedicated mobile skill when the application is bootstrapped.

Do not choose state-management, navigation, persistence, analytics, or other implementation libraries through this contract alone.
