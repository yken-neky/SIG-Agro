#!/bin/bash
# SIG-Agro gRPC Examples
# Run individual commands to test microservices

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================"
echo "  SIG-Agro gRPC Examples"
echo "================================${NC}\n"

# ==================================================
# USER SERVICE EXAMPLES (Port 50051)
# ==================================================

echo -e "${GREEN}1. USER SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Register a new user:"
echo "-------------------"
grpcurl -plaintext \
  -d '{
    "email": "juan.garcia@farm.es",
    "password": "securepass123",
    "full_name": "Juan García López",
    "phone": "+34-555-0123"
  }' \
  localhost:50051 \
  user.UserService/Register

echo ""
echo "Login:"
echo "------"
grpcurl -plaintext \
  -d '{
    "email": "juan.garcia@farm.es",
    "password": "securepass123"
  }' \
  localhost:50051 \
  user.UserService/Login

echo ""
echo "Validate Token (use token from login response):"
echo "-----------------------------------------------"
# Save token from login: TOKEN="<token_from_login>"
read -p "Enter JWT token (or press Enter to skip): " TOKEN

if [ ! -z "$TOKEN" ]; then
  grpcurl -plaintext \
    -d "{\"token\": \"$TOKEN\"}" \
    localhost:50051 \
    user.UserService/ValidateToken
fi

echo ""
echo "Get User Info:"
echo "--------------"
grpcurl -plaintext \
  -d '{"user_id": 1}' \
  localhost:50051 \
  user.UserService/GetUser

echo ""
echo "List Users:"
echo "-----------"
grpcurl -plaintext \
  -d '{"limit": 10, "offset": 0}' \
  localhost:50051 \
  user.UserService/ListUsers

# ==================================================
# PRODUCER SERVICE EXAMPLES (Port 50052)
# ==================================================

echo -e "\n${GREEN}2. PRODUCER SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Create Producer:"
echo "---------------"
grpcurl -plaintext \
  -d '{
    "user_id": 1,
    "name": "Cooperativa García",
    "document_id": "12345678A",
    "phone": "+34-555-0124",
    "email": "info@cooperativa-garcia.es",
    "address": "Carretera Rural km 5, Madrid, España"
  }' \
  localhost:50052 \
  producer.ProducerService/CreateProducer

echo ""
echo "Get Producer:"
echo "-------------"
grpcurl -plaintext \
  -d '{"producer_id": 1}' \
  localhost:50052 \
  producer.ProducerService/GetProducer

echo ""
echo "List Producers:"
echo "---------------"
grpcurl -plaintext \
  -d '{"user_id": 1, "limit": 10, "offset": 0}' \
  localhost:50052 \
  producer.ProducerService/ListProducers

echo ""
echo "Update Producer:"
echo "---------------"
grpcurl -plaintext \
  -d '{
    "producer_id": 1,
    "name": "Cooperativa García Mejorada",
    "phone": "+34-555-0925",
    "email": "contacto@cooperativa-garcia.es",
    "address": "Carretera Rural km 5.5, Madrid, España"
  }' \
  localhost:50052 \
  producer.ProducerService/UpdateProducer

# ==================================================
# PARCEL SERVICE EXAMPLES (Port 50053)
# ==================================================

echo -e "\n${GREEN}3. PARCEL SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Create Parcel with PostGIS Geometry:"
echo "------------------------------------"
echo "(Creating a rectangular parcel near Madrid)"
grpcurl -plaintext \
  -d '{
    "producer_id": 1,
    "name": "Parcela A1 - Maíz",
    "description": "Parcela destinada al cultivo de maíz de primavera",
    "geometry_wkt": "POLYGON((-3.5 40.5, -3.4 40.5, -3.4 40.6, -3.5 40.6, -3.5 40.5))",
    "area_hectares": 5.25,
    "crop_type": "maíz"
  }' \
  localhost:50053 \
  parcel.ParcelService/CreateParcel

echo ""
echo "Create another Parcel (for spatial queries):"
echo "--------------------------------------------"
grpcurl -plaintext \
  -d '{
    "producer_id": 1,
    "name": "Parcela B2 - Trigo",
    "description": "Parcela para cultivo de trigo de invierno",
    "geometry_wkt": "POLYGON((-3.42 40.52, -3.35 40.52, -3.35 40.58, -3.42 40.58, -3.42 40.52))",
    "area_hectares": 8.75,
    "crop_type": "trigo"
  }' \
  localhost:50053 \
  parcel.ParcelService/CreateParcel

echo ""
echo "Get Parcel:"
echo "-----------"
grpcurl -plaintext \
  -d '{"parcel_id": 1}' \
  localhost:50053 \
  parcel.ParcelService/GetParcel

echo ""
echo "List Parcels:"
echo "-------------"
grpcurl -plaintext \
  -d '{"producer_id": 1, "limit": 10, "offset": 0}' \
  localhost:50053 \
  parcel.ParcelService/ListParcels

echo ""
echo "Spatial Query - Find intersecting parcels:"
echo "-------------------------------------------"
echo "(Query for parcels intersecting a polygon)"
grpcurl -plaintext \
  -d '{
    "geometry_wkt": "POLYGON((-3.45 40.52, -3.35 40.52, -3.35 40.58, -3.45 40.58, -3.45 40.52))",
    "query_type": "intersects"
  }' \
  localhost:50053 \
  parcel.ParcelService/QueryByGeometry

# ==================================================
# PRODUCTION SERVICE EXAMPLES (Port 50054)
# ==================================================

echo -e "\n${GREEN}4. PRODUCTION SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Record Production Activity:"
echo "--------------------------"
grpcurl -plaintext \
  -d '{
    "parcel_id": 1,
    "activity_type": "siembra",
    "description": "Siembra de semilla de maíz híbrida",
    "timestamp": '$(date +%s)',
    "metadata": {
      "semilla_variedad": "DK315",
      "cantidad_kg": "25",
      "densidad": "75000_plantas_ha"
    }
  }' \
  localhost:50054 \
  production.ProductionService/RecordActivity

echo ""
echo "Record Fertilization Activity:"
echo "------------------------------"
grpcurl -plaintext \
  -d '{
    "parcel_id": 1,
    "activity_type": "fertilizacion",
    "description": "Aplicación NPK en cobertura",
    "timestamp": '$(date +%s)',
    "metadata": {
      "fertilizante": "NPK_10-10-10",
      "cantidad_kg": "500",
      "metodo": "distribucion_mecanizada"
    }
  }' \
  localhost:50054 \
  production.ProductionService/RecordActivity

echo ""
echo "Record Pest Control Activity:"
echo "-----------------------------"
grpcurl -plaintext \
  -d '{
    "parcel_id": 1,
    "activity_type": "control_plagas",
    "description": "Aplicación de insecticida contra barrenador",
    "timestamp": '$(date +%s)',
    "metadata": {
      "plagas_objetivo": "barrenador_maiz",
      "insecticida": "cipermetrina_25ec",
      "dosis_l_ha": "0.5"
    }
  }' \
  localhost:50054 \
  production.ProductionService/RecordActivity

# ==================================================
# ALERT SERVICE EXAMPLES (Port 50055)
# ==================================================

echo -e "\n${GREEN}5. ALERT SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Create Weather Alert:"
echo "--------------------"
grpcurl -plaintext \
  -d '{
    "parcel_id": 1,
    "alert_type": "clima",
    "severity": "media",
    "message": "Pronóstico de lluvia el próximo fin de semana"
  }' \
  localhost:50055 \
  alert.AlertService/CreateAlert

echo ""
echo "Create Pest Alert:"
echo "-----------------"
grpcurl -plaintext \
  -d '{
    "parcel_id": 1,
    "alert_type": "plagas",
    "severity": "alta",
    "message": "Detección de ácaros rojos en parcela A1"
  }' \
  localhost:50055 \
  alert.AlertService/CreateAlert

echo ""
echo "List Alerts:"
echo "------------"
grpcurl -plaintext \
  -d '{"parcel_id": 1, "severity": "alta", "limit": 10, "offset": 0}' \
  localhost:50055 \
  alert.AlertService/ListAlerts

# ==================================================
# NOTIFICATION SERVICE EXAMPLES (Port 50056)
# ==================================================

echo -e "\n${GREEN}6. NOTIFICATION SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Send Push Notification:"
echo "----------------------"
grpcurl -plaintext \
  -d '{
    "user_id": 1,
    "notification_type": "alerta",
    "channel": "push",
    "title": "Alerta de Plagas",
    "message": "Se ha detectado presencia de ácaros rojos en tu parcela A1",
    "metadata": {
      "parcel_id": "1",
      "alert_id": "1",
      "severity": "alta"
    }
  }' \
  localhost:50056 \
  notification.NotificationService/SendNotification

echo ""
echo "Send Email Notification:"
echo "-----------------------"
grpcurl -plaintext \
  -d '{
    "user_id": 1,
    "notification_type": "reporte",
    "channel": "email",
    "title": "Reporte semanal de actividades",
    "message": "Tu reporte semanal de producción está disponible",
    "metadata": {
      "report_id": "1",
      "periodo": "2024-01-08_a_2024-01-14"
    }
  }' \
  localhost:50056 \
  notification.NotificationService/SendNotification

echo ""
echo "List Notifications:"
echo "------------------"
grpcurl -plaintext \
  -d '{"user_id": 1, "unread_only": true, "limit": 10, "offset": 0}' \
  localhost:50056 \
  notification.NotificationService/ListNotifications

# ==================================================
# REPORT SERVICE EXAMPLES (Port 50057)
# ==================================================

echo -e "\n${GREEN}7. REPORT SERVICE EXAMPLES${NC}"
echo "=================================================="
echo ""

echo "Generate Summary Report:"
echo "-----------------------"
grpcurl -plaintext \
  -d '{
    "producer_id": 1,
    "report_type": "resumen",
    "start_date": '$(date -d "30 days ago" +%s)',
    "end_date": '$(date +%s)',
    "parcel_ids": [1]
  }' \
  localhost:50057 \
  report.ReportService/GenerateReport

echo ""
echo "Generate Detailed Production Report:"
echo "-----------------------------------"
grpcurl -plaintext \
  -d '{
    "producer_id": 1,
    "report_type": "produccion_detallada",
    "start_date": '$(date -d "1 month ago" +%s)',
    "end_date": '$(date +%s)',
    "parcel_ids": [1, 2]
  }' \
  localhost:50057 \
  report.ReportService/GenerateReport

echo ""
echo "List Reports:"
echo "-------------"
grpcurl -plaintext \
  -d '{"producer_id": 1, "limit": 10, "offset": 0}' \
  localhost:50057 \
  report.ReportService/ListReports

echo ""
echo -e "${GREEN}================================"
echo "  Examples completed!"
echo "================================${NC}\n"

echo "Notes:"
echo "------"
echo "1. Replace user_id, producer_id, parcel_id with actual IDs from your database"
echo "2. Use token from Login response with ValidateToken"
echo "3. Timestamps should be Unix epoch seconds"
echo "4. WKT geometries use coordinates as (longitude latitude)"
echo "5. Metadata fields are flexible JSON key-value pairs"
echo ""
echo "Learn more about gRPC:"
echo "  https://grpc.io/docs/languages/go/"
echo ""
