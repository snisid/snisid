# 🗺️ NATIONAL GEOINT PLATFORM

> **Objectif** : Intelligence géospatiale nationale pour décisions terrain intelligentes.

---

## 1. CAPACITÉS

| Fonction | Support |
|----------|:-------:|
| Regional analytics | ✅ |
| Infrastructure mapping | ✅ |
| Deployment heatmaps | ✅ |
| Disaster visualization | ✅ |

---

## 2. STACK TECHNIQUE

| Domaine | Outil |
|---------|-------|
| Base spatiale | PostGIS (PostgreSQL extension) |
| Tuilage / WMS / WFS | GeoServer |
| Cartographie web | OpenLayers |
| Routing | pgRouting / OSRM souverain |
| Imagerie satellite | Sentinel Hub (mirroir local) + drones nationaux |
| Format vectoriel | GeoJSON, MVT (Mapbox Vector Tiles) |
| Format raster | COG (Cloud Optimized GeoTIFF) |
| Catalogue | pycsw (CSW/ISO19115) |

---

## 3. COUCHES NATIONALES

| Couche | Description | Source |
|--------|-------------|--------|
| `admin.departements` | 10 départements | CNIGS |
| `admin.communes` | 145 communes | CNIGS |
| `admin.sections_communales` | 571 SC | CNIGS |
| `infra.routes` | Réseau routier classé | MTPTC |
| `infra.electricite` | Réseau EDH | EDH |
| `infra.telecom` | Antennes télécom | CONATEL |
| `snisid.bureaux` | Bureaux d'enrôlement | SNISID |
| `snisid.agents_deployes` | Position agents temps réel | SNISID mobile |
| `pop.densite` | Densité population (gold) | INSS + SNISID |
| `risk.zones_seismiques` | Aléa sismique | URGéo |
| `risk.zones_inondables` | Aléa inondation | CNIGS |
| `crisis.zones_impactees` | Live | Crisis Engine |

---

## 4. ARCHITECTURE

```
[Frontends Web/Mobile] (OpenLayers + MVT)
        │
        ▼
   ┌─────────────────────┐
   │ GeoServer (WMS/WMTS │
   │ /WFS, MVT, OGC API) │
   └─────────┬───────────┘
             │
   ┌─────────┴───────────┐
   ▼                     ▼
[PostGIS]          [COG Tile Cache MinIO]
   │
   └───── ETL: Spark + GeoPandas → Lakehouse Gold spatial
```

---

## 5. EXEMPLE — DEPLOYMENT HEATMAP

```sql
-- Heatmap agents déployés (PostGIS)
SELECT ST_AsMVT(t, 'agents_heat', 4096, 'geom') AS mvt
FROM (
  SELECT id,
         ST_AsMVTGeom(
           ST_Transform(geom, 3857),
           ST_TileEnvelope({z},{x},{y}),
           4096, 64, true) AS geom,
         intensity
  FROM (
    SELECT ST_SnapToGrid(geom, 0.01) AS geom,
           COUNT(*) AS intensity
    FROM snisid.agents_deployes
    WHERE last_ping > now() - interval '15 minutes'
    GROUP BY ST_SnapToGrid(geom, 0.01)
  ) g
) t;
```

---

## 6. EXEMPLE — DISASTER VISUALIZATION (OpenLayers)

```javascript
import Map from 'ol/Map';
import View from 'ol/View';
import VectorTileLayer from 'ol/layer/VectorTile';
import VectorTileSource from 'ol/source/VectorTile';
import MVT from 'ol/format/MVT';

const impactLayer = new VectorTileLayer({
  source: new VectorTileSource({
    format: new MVT(),
    url: 'https://geoserver.snisid.ht/gwc/service/tms/1.0.0/' +
         'crisis:zones_impactees@EPSG:3857@pbf/{z}/{x}/{-y}.pbf'
  }),
  style: feature => styleByImpact(feature.get('impact_score'))
});

const map = new Map({
  target: 'map',
  layers: [baseLayerHaiti, impactLayer],
  view: new View({ center: [-8092000, 2120000], zoom: 8 })
});
```

---

## 7. ANALYTIQUES SPATIALES

| Analyse | Méthode |
|---------|---------|
| Accessibilité bureau le plus proche | Isochrones pgRouting |
| Zones blanches d'identification | Spatial join densité × bureaux |
| Optimisation tournées mobiles | TSP / VRP avec OR-Tools |
| Détection clusters fraude géo | DBSCAN spatial |
| Couverture biométrique régionale | Choropleth % population enrôlée |

---

## 8. SOUVERAINETÉ GEOINT

- Toutes données cartographiques nationales hébergées localement
- Imagerie satellite mise en miroir sur MinIO souverain
- Pas d'API tierce externe pour cartes opérationnelles
- Géocodage interne (Nominatim local) — pas de Google/Mapbox cloud

---

## 9. CAS D'USAGE TRANSVERSES

| Cas | Bénéficiaire |
|-----|--------------|
| Planification bureaux mobiles | Direction Opérations |
| Évaluation accessibilité services | Présidence |
| Coordination secours | DPC |
| Surveillance frontalière | Sécurité nationale |
| Vote / découpage électoral | CEP |

---

## 10. KPI GEOINT

| KPI | Cible |
|-----|-------|
| Latence tuile MVT (P95) | < 200 ms |
| Disponibilité GeoServer | > 99.9 % |
| Fraîcheur position agents | < 30 s |
| Couverture cartographique nationale | 100 % communes |
