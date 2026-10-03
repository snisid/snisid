# AFIS legacy draft (non-compilable snapshot)

Ce dossier contient une **maquette antérieure** du service AFIS (imports vers un
module `github.com/snisid/afis-svc` inexistant, aucune référence dans `go.work`,
aucun consommateur). Elle a été neutralisée car elle cassait le build global.

La source de vérité du service AFIS est : `services/afis-svc` (module
`github.com/snisid/platform/services/afis-svc`, référencé dans `go.work`,
testé et compilé). Les migrations SQL associées restent dans
`snisid-afis-ht/migrations/afis`.
