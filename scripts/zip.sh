#!/bin/bash
echo "PACKAGING SNISID v1.0 PLATFORM..."
zip -r sniseid-platform.zip . -x "*.git*" -x "node_modules/*" -x "*.exe" -x "*.zip"
echo "PACKAGE CREATED: sniseid-platform.zip"
