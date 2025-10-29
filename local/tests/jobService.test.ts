import { expect } from "@std/expect";
import axios from "axios";
import MockAdapter from "axios-mock-adapter";
import { Job, JobService, Response } from "../jobService.ts";
let mockAxios: MockAdapter;
let jobService: JobService;

Deno.test.beforeEach(() => {
  // Remove the ignore when https://github.com/ctimmerm/axios-mock-adapter/issues/400 is fixed
  // deno-lint-ignore no-explicit-any
  mockAxios = new MockAdapter(axios as any);
  jobService = new JobService("http://mock-endpoint.com");
});

Deno.test.afterEach(() => {
  mockAxios.reset();
});

// Using explicit Deno.test instead of t.step so the before/after hooks runs for each test
Deno.test("should send jobs and return the response", async () => {
  const jobs: Job[] = [
    {
      title: "Software Engineer",
      company: "Company1",
      keyword: "Go",
      link: "http://example.com/job1",
    },
    {
      title: "Data Analyst",
      company: "Company2",
      keyword: "Python",
      link: "http://example.com/job2",
    },
    {
      title: "Financial Advisor",
      company: "Company3",
      keyword: "Finance",
      link: "http://example.com/job3",
    },
  ];

  const expectedResponse: Response = { total: 3, uncached: 1, duplicates: 2 };
  mockAxios.onPost("http://mock-endpoint.com").reply(200, expectedResponse);
  const response = await jobService.sendJobs(jobs);
  expect(response).toEqual(expectedResponse);
});

Deno.test("should throw an error if the request fails", async () => {
  const jobs: Job[] = [
    {
      title: "Software Engineer",
      company: "Company1",
      keyword: "Go",
      link: "http://example.com/job1",
    },
  ];

  mockAxios.onPost("http://mock-endpoint.com").networkError();
  await expect(jobService.sendJobs(jobs)).rejects.toThrow("Network Error");
});
