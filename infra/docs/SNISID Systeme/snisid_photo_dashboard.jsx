import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { useState } from "react";

// Lightweight inline icons (replaces lucide-react to avoid CDN/module issues)
const CameraIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M4 7h3l2-2h6l2 2h3v14H4V7z" strokeWidth="2" />
    <circle cx="12" cy="14" r="3" strokeWidth="2" />
  </svg>
);

const ScanIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M4 7V4h3M20 7V4h-3M4 17v3h3M20 17v3h-3" strokeWidth="2" />
    <circle cx="12" cy="12" r="3" strokeWidth="2" />
  </svg>
);

const AlertIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M12 9v4m0 4h.01M10.29 3.86l-8.3 14.37A1 1 0 0 0 3 20h18a1 1 0 0 0 .86-1.5L13.71 3.86a1 1 0 0 0-1.72 0z" strokeWidth="2" />
  </svg>
);

const CheckIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M20 6L9 17l-5-5" strokeWidth="2" />
  </svg>
);

export default function SNISIDDashboard() {
  const [file, setFile] = useState(null);
  const [uploadStatus, setUploadStatus] = useState("idle");

  const handleUpload = () => {
    if (!file) return;
    setUploadStatus("processing");

    setTimeout(() => {
      setUploadStatus("done");
    }, 1200);
  };

  return (
    <div className="p-6 grid gap-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">SNISID Photo Intelligence Dashboard</h1>
        <Badge variant="outline">SOC Live</Badge>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <CameraIcon />
              <h2 className="font-semibold">Photo Intake</h2>
            </div>
            <Input
              type="file"
              className="mt-3"
              onChange={(e) => setFile(e.target.files?.[0] || null)}
            />
            <Button className="mt-3 w-full" onClick={handleUpload}>
              {uploadStatus === "processing" ? "Processing..." : "Upload for Analysis"}
            </Button>
            {uploadStatus === "done" && (
              <div className="mt-2 flex items-center gap-2 text-green-600">
                <CheckIcon /> Analysis complete
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <ScanIcon />
              <h2 className="font-semibold">Face Verification</h2>
            </div>
            <div className="mt-3 space-y-2">
              <div className="flex justify-between">
                <span>Identity Match</span>
                <Badge>98.4%</Badge>
              </div>
              <div className="flex justify-between">
                <span>Liveness Score</span>
                <Badge>96.1%</Badge>
              </div>
              <div className="flex justify-between">
                <span>Deepfake Risk</span>
                <Badge variant="destructive">Low</Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <AlertIcon />
              <h2 className="font-semibold">Threat Alerts</h2>
            </div>
            <div className="mt-3 space-y-2 text-sm">
              <div className="flex justify-between">
                <span>Face Spoof Attempt</span>
                <Badge variant="destructive">Blocked</Badge>
              </div>
              <div className="flex justify-between">
                <span>Duplicate Identity</span>
                <Badge variant="outline">Investigating</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="live">
        <TabsList>
          <TabsTrigger value="live">Live Feed</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="audit">Audit Logs</TabsTrigger>
        </TabsList>

        <TabsContent value="live" className="mt-4">
          <Card>
            <CardContent className="p-4 text-sm">
              Live SNISID verification stream active...
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics" className="mt-4">
          <Card>
            <CardContent className="p-4">
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Total Verifications</span>
                  <span>12,482</span>
                </div>
                <div className="flex justify-between">
                  <span>Fraud Attempts</span>
                  <span>214</span>
                </div>
                <div className="flex justify-between">
                  <span>Deepfake Detected</span>
                  <span>37</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="audit" className="mt-4">
          <Card>
            <CardContent className="p-4 text-sm">
              Immutable audit logs (Kafka-backed) for all identity events.
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Basic sanity test cases (non-executed reference) */}
      {/*
        TEST CASES:
        1. Upload file → uploadStatus becomes "processing"
        2. After 1.2s → uploadStatus becomes "done"
        3. No file selected → handleUpload does nothing
      */}
    </div>
  );
}
