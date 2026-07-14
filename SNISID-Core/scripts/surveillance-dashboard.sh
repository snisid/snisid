#!/bin/bash
###############################################################################
# SNISID REAL-TIME SURVEILLANCE DASHBOARD
# Inspiré du système "NSA XKeyscore" et "Chine Skynet/Safe City"
# Objectif: Surveillance temps réel des flux de données, alertes automatiques
###############################################################################

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/../config/dashboard.conf"
LOG_FILE="${SCRIPT_DIR}/logs/surveillance.log"
ALERT_FILE="${SCRIPT_DIR}/alerts/active_alerts.json"
DASHBOARD_PORT=3000

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log() {
    local level=$1
    shift
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $*" | tee -a "$LOG_FILE"
}

log_info() { log "${BLUE}INFO${NC}" "$@"; }
log_success() { log "${GREEN}SUCCESS${NC}" "$@"; }
log_warning() { log "${YELLOW}WARNING${NC}" "$@"; }
log_error() { log "${RED}ERROR${NC}" "$@"; }
log_alert() { log "${RED}⚠️ ALERTE${NC}" "$@"; }

# Créer les répertoires nécessaires
setup_directories() {
    mkdir -p "${SCRIPT_DIR}/logs" "${SCRIPT_DIR}/alerts" "${SCRIPT_DIR}/reports" "${SCRIPT_DIR}/../data/surveillance"
    log_info "Répertoires de surveillance créés."
}

# Installation des composants du dashboard
install_dashboard_components() {
    log_info "Installation des composants du dashboard de surveillance..."
    
    # Vérifier si Node.js est installé
    if ! command -v node &> /dev/null; then
        log_warning "Node.js non détecté. Installation en cours..."
        curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
        apt-get install -y nodejs
    fi
    
    # Installer PM2 pour la gestion des processus
    if ! command -v pm2 &> /dev/null; then
        npm install -g pm2
    fi
    
    log_success "Composants du dashboard installés."
}

# Configuration de la base de données Time-Series (InfluxDB)
setup_timeseries_db() {
    log_info "Configuration de la base de données Time-Series pour les métriques..."
    
    cat > "${SCRIPT_DIR}/../config/influxdb-config.toml" <<EOF
# Configuration InfluxDB 2.x pour SNISID
# Stockage des métriques de surveillance (Style NSA/CNSA)

[http]
  bind-address = ":8086"
  https-enabled = true
  https-certificate = "/etc/ssl/snisid/influxdb.crt"
  https-private-key = "/etc/ssl/snisid/influxdb.key"

[meta]
  dir = "/var/lib/influxdb2/meta"
  retention-autocreate = true

[data]
  dir = "/var/lib/influxdb2/data"
  wal-dir = "/var/lib/influxdb2/wal"
  max-concurrent-queries = 100
  
[coordinator]
  max-select-point = 0
  max-select-series = 0
  
[retention-policies]
  # Rétention des données selon standards MLPS Chine
  [[retention-policies.autogen]]
    name = "autogen"
    duration = "8760h"  # 1 an
    shard-group-duration = "168h"
    replica-n = 1
    default = true
    
  [[retention-policies.long_term]]
    name = "long_term"
    duration = "87600h"  # 10 ans (Archives légales)
    shard-group-duration = "720h"
    replica-n = 3
EOF

    log_success "Configuration InfluxDB créée."
}

# Création des règles d'alerte intelligentes
create_alert_rules() {
    log_info "Création des règles d'alerte intelligentes (Style FBI Sentinel)..."
    
    cat > "${SCRIPT_DIR}/../config/alert-rules.json" <<EOF
{
  "version": "2.0",
  "agency": "SNISID-Haiti",
  "classification": "TOP SECRET",
  "alert_rules": [
    {
      "id": "RULE-001",
      "name": "Détection de mouvement suspect aux frontières",
      "description": "Alerte lorsqu'un individu surveillé traverse une frontière sans autorisation",
      "severity": "CRITICAL",
      "trigger": {
        "type": "geo_fence_breach",
        "conditions": {
          "person_risk_level": ["HIGH", "CRITICAL"],
          "location_type": ["border_crossing", "airport", "seaport"],
          "time_window": "real_time"
        }
      },
      "actions": [
        {"type": "sms", "recipients": ["+509-XXX-XXXX", "+509-XXX-XXXX"]},
        {"type": "email", "recipients": ["ops@snisid.gouv.ht"]},
        {"type": "dashboard_popup", "priority": 1},
        {"type": "auto_create_case", "system": "DEA-like"}
      ]
    },
    {
      "id": "RULE-002",
      "name": "Transaction financière anormale",
      "description": "Détection de blanchiment d'argent ou financement terroriste",
      "severity": "HIGH",
      "trigger": {
        "type": "financial_anomaly",
        "conditions": {
          "amount_threshold": 50000,
          "currency": ["USD", "HTG", "EUR"],
          "pattern": ["structuring", "rapid_movement", "shell_company"],
          "linked_entities": ["sanctioned_individual", "pep"]
        }
      },
      "actions": [
        {"type": "alert_queue", "destination": "financial_crimes_unit"},
        {"type": "generate_report", "format": "PDF"},
        {"type": "freeze_assets", "auto": false, "requires_approval": true}
      ]
    },
    {
      "id": "RULE-003",
      "name": "Communication cryptée suspecte",
      "description": "Détection d'utilisation d'applications de communication cryptée par des cibles",
      "severity": "MEDIUM",
      "trigger": {
        "type": "comms_pattern",
        "conditions": {
          "app_type": ["signal", "telegram", "whatsapp", "threema"],
          "frequency": ">50_messages_per_hour",
          "participants_count": ">5",
          "time_of_day": ["night", "early_morning"]
        }
      },
      "actions": [
        {"type": "log_enhanced_monitoring"},
        {"type": "flag_for_analysis"},
        {"type": "correlate_with_other_intel"}
      ]
    },
    {
      "id": "RULE-004",
      "name": "Regroupement de gangs détecté",
      "description": "Algorithmes de détection de rassemblements suspects (Style China Skynet)",
      "severity": "HIGH",
      "trigger": {
        "type": "crowd_anomaly",
        "conditions": {
          "min_people": 10,
          "known_gang_members": ">=3",
          "location_type": ["residential", "warehouse", "remote_area"],
          "time_window": "after_22h"
        }
      },
      "actions": [
        {"type": "activate_nearby_cctv"},
        {"type": "dispatch_quick_response_team"},
        {"type": "record_evidence"}
      ]
    },
    {
      "id": "RULE-005",
      "name": "Tentative d'accès non autorisé au système",
      "description": "Détection d'intrusion cybernétique (Style NSA Cybersecurity)",
      "severity": "CRITICAL",
      "trigger": {
        "type": "cyber_intrusion",
        "conditions": {
          "failed_login_attempts": ">5",
          "source_ip_reputation": "malicious",
          "geo_location": "outside_haiti",
          "attack_signature": ["brute_force", "sql_injection", "zero_day"]
        }
      },
      "actions": [
        {"type": "block_ip_immediately"},
        {"type": "isolate_affected_system"},
        {"type": "notify_cyber_command"},
        {"type": "initiate_forensic_capture"}
      ]
    }
  ],
  "escalation_matrix": {
    "CRITICAL": {
      "notification_delay": "0_minutes",
      "recipients": ["Directeur_General", "Ministre_Interieur", "President"],
      "auto_response": true
    },
    "HIGH": {
      "notification_delay": "5_minutes",
      "recipients": ["Directeur_Operations", "Chef_Antigang"],
      "auto_response": false
    },
    "MEDIUM": {
      "notification_delay": "30_minutes",
      "recipients": ["Analyste_Senior"],
      "auto_response": false
    },
    "LOW": {
      "notification_delay": "24_hours",
      "recipients": ["Analyste_Junior"],
      "auto_response": false
    }
  }
}
EOF

    log_success "Règles d'alerte intelligentes créées (5 règles actives)."
}

# Déploiement du dashboard React/Node.js
deploy_dashboard_ui() {
    log_info "Déploiement de l'interface de surveillance..."
    
    local dashboard_dir="${SCRIPT_DIR}/../dashboard-ui"
    
    if [ ! -d "$dashboard_dir" ]; then
        mkdir -p "$dashboard_dir"
        cd "$dashboard_dir"
        
        # Initialiser le projet React
        npm init -y
        npm install react react-dom express socket.io axios chart.js react-chartjs-2 moment
        
        # Créer le serveur Express
        cat > server.js <<'SERVEREOF'
const express = require('express');
const http = require('http');
const socketIO = require('socket.io');
const path = require('path');

const app = express();
const server = http.createServer(app);
const io = socketIO(server, {
    cors: { origin: "*", methods: ["GET", "POST"] }
});

const PORT = process.env.PORT || 3000;

// Servir les fichiers statiques
app.use(express.static(path.join(__dirname, 'public')));

// API Endpoints
app.get('/api/alerts', (req, res) => {
    // Récupérer les alertes depuis Elasticsearch/InfluxDB
    res.json({ alerts: [], timestamp: new Date() });
});

app.get('/api/metrics', (req, res) => {
    // Récupérer les métriques de surveillance
    res.json({ 
        active_surveillance: 0,
        total_persons_of_interest: 0,
        alerts_today: 0 
    });
});

// WebSocket pour les mises à jour en temps réel
io.on('connection', (socket) => {
    console.log('Client connecté au dashboard:', socket.id);
    
    socket.on('subscribe', (channel) => {
        socket.join(channel);
        console.log(`Client ${socket.id} subscribed to ${channel}`);
    });
    
    socket.on('disconnect', () => {
        console.log('Client déconnecté:', socket.id);
    });
});

server.listen(PORT, '0.0.0.0', () => {
    console.log(`SNISID Dashboard running on port ${PORT}`);
    console.log(`Access: http://localhost:${PORT}`);
});
SERVEREOF

        # Créer le dossier public
        mkdir -p public
        
        # Créer le fichier HTML principal
        cat > public/index.html <<'HTMLEOF'
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SNISID - Centre de Commandement National</title>
    <script src="https://cdn.socket.io/4.5.4/socket.io.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: #0a0a0a; 
            color: #00ff00;
            overflow-x: hidden;
        }
        .header {
            background: linear-gradient(90deg, #1a1a1a 0%, #0d0d0d 100%);
            padding: 20px;
            border-bottom: 3px solid #00ff00;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .header h1 { font-size: 24px; text-transform: uppercase; letter-spacing: 2px; }
        .status-badge { 
            background: #00ff00; 
            color: #000; 
            padding: 5px 15px; 
            border-radius: 20px; 
            font-weight: bold;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0% { opacity: 1; }
            50% { opacity: 0.7; }
            100% { opacity: 1; }
        }
        .dashboard-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 20px;
            padding: 20px;
        }
        .panel {
            background: #111;
            border: 1px solid #333;
            border-radius: 8px;
            padding: 20px;
            min-height: 300px;
        }
        .panel h2 { 
            font-size: 16px; 
            margin-bottom: 15px; 
            border-bottom: 1px solid #333;
            padding-bottom: 10px;
            color: #fff;
        }
        .alert-item {
            background: #1a0a0a;
            border-left: 4px solid #ff0000;
            padding: 15px;
            margin-bottom: 10px;
            border-radius: 4px;
        }
        .alert-item.critical { border-left-color: #ff0000; background: #2a0a0a; }
        .alert-item.high { border-left-color: #ff6600; }
        .alert-item.medium { border-left-color: #ffcc00; }
        .metric-value { font-size: 48px; font-weight: bold; color: #00ff00; }
        .metric-label { font-size: 14px; color: #888; margin-top: 5px; }
        .map-container { height: 400px; background: #0a0a0a; border-radius: 8px; }
        .footer {
            text-align: center;
            padding: 20px;
            color: #666;
            font-size: 12px;
            border-top: 1px solid #333;
        }
        .classification-banner {
            background: #ff0000;
            color: white;
            text-align: center;
            padding: 5px;
            font-weight: bold;
            text-transform: uppercase;
        }
    </style>
</head>
<body>
    <div class="classification-banner">TOP SECRET // NOFORN // Gouvernement Haïtien</div>
    
    <div class="header">
        <div>
            <h1>🇭🇹 SNISID - Centre de Commandement National</h1>
            <small>Système National d'Identification Sécurisée</small>
        </div>
        <div class="status-badge">● SYSTÈME OPÉRATIONNEL</div>
    </div>
    
    <div class="dashboard-grid">
        <!-- Panel 1: Métriques Globales -->
        <div class="panel">
            <h2>📊 MÉTRIQUES EN TEMPS RÉEL</h2>
            <div class="metric-value" id="active-surveillance">0</div>
            <div class="metric-label">Surveillances Actives</div>
            <hr style="border-color: #333; margin: 20px 0;">
            <div class="metric-value" id="total-poi">0</div>
            <div class="metric-label">Personnes d'Intérêt</div>
            <hr style="border-color: #333; margin: 20px 0;">
            <div class="metric-value" id="alerts-today" style="color: #ff0000;">0</div>
            <div class="metric-label">Alertes Aujourd'hui</div>
        </div>
        
        <!-- Panel 2: Alertes Critiques -->
        <div class="panel">
            <h2>🚨 ALERTES CRITIQUES</h2>
            <div id="alerts-container">
                <div class="alert-item critical">
                    <strong>RULE-001:</strong> Mouvement suspect détecté à la frontière nord
                    <br><small>Il y a 2 minutes - En attente de validation</small>
                </div>
                <div class="alert-item high">
                    <strong>RULE-002:</strong> Transaction financière anormale: $75,000 USD
                    <br><small>Il y a 15 minutes - Unité crimes financiers notifiée</small>
                </div>
            </div>
        </div>
        
        <!-- Panel 3: Carte Interactive -->
        <div class="panel">
            <h2>🗺️ CARTE DE SURVEILLANCE NATIONALE</h2>
            <div class="map-container" id="map">
                <!-- Intégration Leaflet/OpenStreetMap ici -->
                <div style="height: 100%; display: flex; align-items: center; justify-content: center; color: #666;">
                    [Carte interactive en chargement...]
                </div>
            </div>
        </div>
        
        <!-- Panel 4: Graphique d'Activité -->
        <div class="panel">
            <h2>📈 ACTIVITÉ DES 24 DERNIÈRES HEURES</h2>
            <canvas id="activityChart"></canvas>
        </div>
    </div>
    
    <div class="footer">
        SNISID v2.0 © 2024 - République d'Haïti<br>
        Classification: TOP SECRET // NOFORN<br>
        Inspiré des architectures NSA DataWave, FBI Palantir, Chine Skynet
    </div>
    
    <script>
        // Connexion WebSocket
        const socket = io('http://localhost:3000');
        
        socket.on('connect', () => {
            console.log('Connecté au serveur SNISID');
            socket.emit('subscribe', 'alerts');
            socket.emit('subscribe', 'metrics');
        });
        
        socket.on('new_alert', (data) => {
            console.log('Nouvelle alerte:', data);
            // Mettre à jour l'UI
            document.getElementById('alerts-today').innerText = 
                parseInt(document.getElementById('alerts-today').innerText) + 1;
        });
        
        // Initialiser le graphique
        const ctx = document.getElementById('activityChart').getContext('2d');
        const activityChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'],
                datasets: [{
                    label: 'Événements détectés',
                    data: [12, 8, 25, 45, 38, 52],
                    borderColor: '#00ff00',
                    backgroundColor: 'rgba(0, 255, 0, 0.1)',
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                scales: {
                    y: { beginAtZero: true, grid: { color: '#333' } },
                    x: { grid: { color: '#333' } }
                },
                plugins: {
                    legend: { labels: { color: '#fff' } }
                }
            }
        });
        
        // Simulation de données en temps réel
        setInterval(() => {
            document.getElementById('active-surveillance').innerText = 
                Math.floor(Math.random() * 50) + 100;
            document.getElementById('total-poi').innerText = 
                Math.floor(Math.random() * 200) + 500;
        }, 5000);
    </script>
</body>
</html>
HTMLEOF

        log_success "Interface dashboard déployée avec succès."
    else
        log_info "Le dashboard existe déjà. Skipping..."
    fi
}

# Démarrer tous les services
start_services() {
    log_info "Démarrage des services de surveillance..."
    
    # Démarrer le serveur dashboard avec PM2
    cd "${SCRIPT_DIR}/../dashboard-ui"
    pm2 start server.js --name snisid-dashboard --port $DASHBOARD_PORT
    
    log_success "Dashboard démarré sur le port ${DASHBOARD_PORT}"
    log_info "Accès: http://localhost:${DASHBOARD_PORT}"
    log_warning "IMPORTANT: Restreindre l'accès au réseau sécurisé uniquement!"
}

# Main execution
main() {
    log_info "=========================================="
    log_info "SNISID Real-Time Surveillance System"
    log_info "Inspiré: NSA XKeyscore / Chine Skynet / FBI Sentinel"
    log_info "=========================================="
    
    setup_directories
    install_dashboard_components
    setup_timeseries_db
    create_alert_rules
    deploy_dashboard_ui
    start_services
    
    log_success "=========================================="
    log_success "SYSTÈME DE SURVEILLANCE OPÉRATIONNEL"
    log_success "Dashboard: http://localhost:${DASHBOARD_PORT}"
    log_success "Logs: ${LOG_FILE}"
    log_success "Alertes: ${ALERT_FILE}"
    log_success "=========================================="
    
    echo ""
    echo "🇭🇹 PROCHAINES ÉTAPES:"
    echo "1. Configurer les flux de données (télécoms, CCTV, frontières)"
    echo "2. Former les analystes du Centre National de Sécurité"
    echo "3. Effectuer un test de simulation de crise"
    echo "4. Obtenir la certification du Ministre de l'Intérieur"
}

main "$@"
