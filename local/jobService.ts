import axios from "axios";

export interface Job {
  title: string;
  company: string;
  keyword: string;
  link: string;
}

export interface Response {
  total: number;
  uncached: number;
  duplicates: number;
}

export class JobService {
  private endpoint: string;

  constructor(endpoint: string) {
    this.endpoint = endpoint;
  }

  async sendJobs(jobs: Job[]): Promise<Response> {
    const response = await axios.post(this.endpoint, { jobs });
    return response.data as Response;
  }
}
