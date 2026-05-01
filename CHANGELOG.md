# Changelog

## [3.0.0] - 2026-05-01

### Changed

- **Per-user traffic metrics are now exposed.** Previously, `user>>>` entries from the Xray Stats API were intentionally skipped to avoid cardinality concerns. This fork removes that restriction.

  The following metrics now include `dimension="user"` with `target` set to the user's email/name as configured in 3x-ui or Xray:

  ```
  xray_traffic_uplink_bytes_total{dimension="user", target="alice-phone"}
  xray_traffic_downlink_bytes_total{dimension="user", target="alice-phone"}
  ```

  This enables per-user dashboards in Grafana — top consumers, individual traffic history, and "who is active now" panels.

### Notes

- User stats require the `email` field to be set for each client in 3x-ui (or the equivalent in bare Xray config). Without it, Xray does not generate per-user stats entries.
- `policy.levels.0.statsUserUplink` and `statsUserDownlink` must be `true` in the Xray config.
- Stats are created lazily — a user entry only appears after the first connection.
