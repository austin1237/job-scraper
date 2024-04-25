import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import { JobService, Job, Response } from '../src/jobService';

describe('JobService', () => {
    let mockAxios: MockAdapter;
    let jobService: JobService;

    beforeEach(() => {
        mockAxios = new MockAdapter(axios);
        jobService = new JobService('http://mock-endpoint.com');
    });

    afterEach(() => {
        mockAxios.reset();
    });

    it('should send jobs and return the response', async () => {
        const jobs: Job[] = [
            { title: 'Software Engineer', company: 'Company1', keyword: 'Go', link: 'http://example.com/job1' },
            { title: 'Data Analyst', company: 'Company2', keyword: 'Python', link: 'http://example.com/job2' },
            { title: 'Financial Advisor', company: 'Company3', keyword: 'Finance', link: 'http://example.com/job3' },
        ];

        const expectedResponse: Response = { total: 3, uncached: 1, duplicates: 2 };

        mockAxios.onPost('http://mock-endpoint.com').reply(200, expectedResponse);

        const response = await jobService.sendJobs(jobs);

        expect(response).toEqual(expectedResponse);
    });

    it('should throw an error if the request fails', async () => {
        const jobs: Job[] = [
            { title: 'Software Engineer', company: 'Company1', keyword: 'Go', link: 'http://example.com/job1' },
        ];

        mockAxios.onPost('http://mock-endpoint.com').networkError();

        await expect(jobService.sendJobs(jobs)).rejects.toThrow('Network Error');
    });
});
