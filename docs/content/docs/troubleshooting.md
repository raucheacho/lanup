---
title: "Troubleshooting"
weight: 6
---

# Troubleshooting

Solutions to common problems you might encounter.

## No Network Interface Detected

**Problem:** lanup can't detect your local IP address.

**Error message:**
```
Error: Failed to detect local IP address
```

**Solutions:**

1. **Check network connection**
   ```bash
   # macOS/Linux
   ifconfig
   
   # Windows
   ipconfig
   ```

2. **Run diagnostics**
   ```bash
   lanup doctor
   ```

3. **Verify network interface is active**
   - Ensure you're connected to Wi-Fi or Ethernet
   - Check that your network adapter is enabled

4. **Check for VPN interference**
   - Some VPNs can interfere with local network detection
   - Try disconnecting from VPN temporarily

---

## Docker Auto-Detection Not Working

**Problem:** Docker containers are not being detected.

**Solutions:**

1. **Verify Docker is running**
   ```bash
   docker ps
   ```

2. **Check Docker permissions**
   ```bash
   # Linux: Add user to docker group
   sudo usermod -aG docker $USER
   # Then log out and back in
   ```

3. **Run diagnostics**
   ```bash
   lanup doctor
   ```

4. **Disable auto-detection and configure manually**
   ```yaml
   # .lanup.yaml
   vars:
     MY_SERVICE_PORT: "http://localhost:8080"
   
   auto_detect:
     docker: false
   ```

---

## Environment File Not Updating

**Problem:** The `.env.local` file is not being updated.

**Solutions:**

1. **Check file permissions**
   ```bash
   ls -la .env.local
   chmod 644 .env.local
   ```

2. **Verify output path**
   ```yaml
   # .lanup.yaml
   output: ".env.local"  # Check this path
   ```

3. **Try dry-run mode**
   ```bash
   lanup start --dry-run
   ```

4. **Check logs**
   ```bash
   lanup logs --tail 50
   ```

5. **Verify disk space**
   ```bash
   df -h
   ```

---

## Network Changes Not Detected in Watch Mode

**Problem:** Watch mode doesn't update when switching networks.

**Solutions:**

1. **Check interval setting**
   ```yaml
   # ~/.lanup/config.yaml
   check_interval: 5  # Increase if needed
   ```

2. **Verify stable connection**
   - Ensure network connection is stable
   - Wait a few seconds after switching networks

3. **Restart watch mode**
   ```bash
   # Stop with Ctrl+C
   lanup start --watch
   ```

4. **Check logs**
   ```bash
   lanup logs --follow
   ```

---

## Permission Denied Errors

**Problem:** lanup can't write files or access logs.

**Error message:**
```
Error: Permission denied
```

**Solutions:**

1. **Check project directory permissions**
   ```bash
   ls -la
   chmod 755 .
   ```

2. **Check ~/.lanup/ permissions**
   ```bash
   ls -la ~/.lanup/
   chmod 755 ~/.lanup/
   chmod 600 ~/.lanup/config.yaml
   ```

3. **Run without sudo**
   - lanup should not require sudo
   - If you used sudo before, fix permissions:
   ```bash
   sudo chown -R $USER:$USER ~/.lanup/
   ```

---

## Services Not Accessible from Other Devices

**Problem:** Other devices can't access the exposed URLs.

**Solutions:**

1. **Verify same network**
   - Ensure all devices are on the same Wi-Fi/network
   - Check device IP addresses are in same subnet

2. **Check firewall settings**
   ```bash
   # macOS: System Preferences > Security & Privacy > Firewall
   # Linux: Check iptables or ufw
   # Windows: Windows Defender Firewall
   ```

3. **Test from development machine**
   ```bash
   curl http://192.168.1.100:3000
   ```

4. **Verify service is running**
   ```bash
   curl http://localhost:3000
   ```

5. **Check network restrictions**
   - Corporate networks may block device-to-device communication
   - Public Wi-Fi often isolates devices
   - Try a different network or use mobile hotspot

---

## Supabase Not Detected

**Problem:** Supabase services are not being detected.

**Solutions:**

1. **Verify Supabase is running**
   ```bash
   supabase status
   ```

2. **Check Supabase CLI installation**
   ```bash
   supabase --version
   ```

3. **Start Supabase**
   ```bash
   supabase start
   ```

4. **Disable auto-detection and configure manually**
   ```yaml
   # .lanup.yaml
   vars:
     SUPABASE_URL: "http://localhost:54321"
   
   auto_detect:
     supabase: false
   ```

---

## Configuration File Not Found

**Problem:** lanup can't find `.lanup.yaml`.

**Error message:**
```
Error: Failed to load project configuration
```

**Solutions:**

1. **Initialize lanup**
   ```bash
   lanup init
   ```

2. **Check current directory**
   ```bash
   pwd
   ls -la .lanup.yaml
   ```

3. **Verify file name**
   - Must be `.lanup.yaml` (with leading dot)
   - Check for typos

---

## Invalid YAML Syntax

**Problem:** Configuration file has syntax errors.

**Error message:**
```
Error: failed to parse config file
```

**Solutions:**

1. **Validate YAML syntax**
   - Use online YAML validator
   - Check indentation (use spaces, not tabs)

2. **Common issues:**
   ```yaml
   # ❌ Wrong - missing quotes
   vars:
     URL: http://localhost:8000
   
   # ✓ Correct
   vars:
     URL: "http://localhost:8000"
   ```

3. **Regenerate config**
   ```bash
   lanup init --force
   ```

---

## Logs Not Being Created

**Problem:** Log file is not being created or updated.

**Solutions:**

1. **Check log path**
   ```yaml
   # ~/.lanup/config.yaml
   log_path: "~/.lanup/logs/lanup.log"
   ```

2. **Verify directory exists**
   ```bash
   ls -la ~/.lanup/logs/
   mkdir -p ~/.lanup/logs/
   ```

3. **Check permissions**
   ```bash
   chmod 755 ~/.lanup/logs/
   ```

4. **Enable logging**
   ```bash
   lanup start --log
   ```

---

## Getting Help

If you're still experiencing issues:

1. **Run diagnostics**
   ```bash
   lanup doctor
   ```

2. **Check logs**
   ```bash
   lanup logs --tail 100
   ```

3. **Enable verbose mode**
   ```bash
   lanup start --verbose
   ```

4. **Report an issue**
   - Visit [GitHub Issues](https://github.com/raucheacho/lanup/issues)
   - Include output from `lanup doctor`
   - Include relevant log entries
   - Describe your environment (OS, network setup)
