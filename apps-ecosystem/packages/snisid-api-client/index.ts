// API Client for SNISID
export class ApiClient {
  constructor(private baseUrl: string) {}

  async get(endpoint: string) {
    console.log(`GET ${this.baseUrl}${endpoint}`);
    return { data: "success" };
  }
}
