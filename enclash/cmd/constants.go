package cmd

const WEB_STATIC_FOLDER = "clash-dashboard"
const CONF_TUN_DNS = `
tun:
  enable: true
  stack: gvisor
  dns-hijack:
    - any:53 
  auto-route: true   
  auto-detect-interface: true 
dns:
  enable: true
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16 
  nameserver:
    - 114.114.114.114 
    - 8.8.8.8
`
const CONF_WEB_UI = `
external-ui: clash-dashboard
`
