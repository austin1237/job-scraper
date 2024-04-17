import { evaluateJobInterest } from '../src/interest';

describe('evaluateJobInterest', () => {
  it('should return null if any of the parameters are null', () => {
    expect(evaluateJobInterest(null, 'Company', 'Job Description')).toBeNull();
    expect(evaluateJobInterest('Software Engineer', null, 'Job Description')).toBeNull();
    expect(evaluateJobInterest('Software Engineer', 'Company', null)).toBeNull();
  });

  it('should return null if no included keywords are found', () => {
    expect(evaluateJobInterest('Software Engineer', 'Company', 'Job Description')).toBeNull();
  });

  it('should return the matched keyword if an included keyword is found and no excluded keywords are found', () => {
    expect(evaluateJobInterest('Software Engineer', 'Company', 'go')).toBe('go');
  });

  it('should return the null if an included keyword is part of another word', () => {
    expect(evaluateJobInterest('Software Engineer', 'Company', 'chicago')).toBeNull();
  });

  it('should return null if an included keyword and an excluded keyword are found', () => {
    expect(evaluateJobInterest('Software Engineer', 'Company', 'node contract')).toBeNull();
  });

  it('should return null if the job title does not match the include regex', () => {
    expect(evaluateJobInterest('Bad Job Title', 'Company', 'node')).toBeNull();
  });

  it('should return null if the job title matches the exclude regex', () => {
    expect(evaluateJobInterest('Front-End Developer', 'Company', 'node')).toBeNull();
  });

  it('should return null if the company name matches the exclude company regex', () => {
    expect(evaluateJobInterest('Software Engineer', 'Consulting', 'node')).toBeNull();
  });
});