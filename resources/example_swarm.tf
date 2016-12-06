resource "ciscodocker_swarm" "swarm" {
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.100:2377"
  force_new_cluster = false
  auto_lock_managers = false
  task_history_retention_limit = 10
  snapshot_interval = 10000
  log_entries_for_slow_followers = 500
  election_tick = 3
  heartbeat_tick = 1
  heartbeat_period = 5000000000
  node_cert_expiry = 7776000000000000
  force_leave = true//to destroy swarm manager
  rotate_manager_token = false//to update swarm
  rotate_worker_token = false//to update swarm
  //manager_token = will compute / string
  //worker_token = will compute / string
  //version = will compute / int

}
