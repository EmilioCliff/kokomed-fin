import { faker } from '@faker-js/faker';
import {
  loanSchema,
  Loan,
  User,
  userSchema,
  branchSchema,
  Branch,
  productSchema,
  Product,
  InactiveLoan,
  Client,
  inactiveLoanSchema,
  clientSchema,
} from '../data/schema';

const formatDate = (date: Date) => date.toISOString().split('T')[0];

function getRandomEnum<T extends string>(enumValues: T[]): T {
  return enumValues[Math.floor(Math.random() * enumValues.length)];
}

function getRandomBoolean(): boolean {
  return Math.random() < 0.5;
}

function getRandomUser(): User {
  return userSchema.parse({
    id: faker.number.int({ min: 1, max: 100 }),
    fullName: faker.person.fullName(),
    phoneNumber: faker.phone.number({ style: 'international' }),
    email: faker.internet.email(),
    role: getRandomEnum(['ADMIN', 'AGENT'] as const),
    branchName: faker.company.name(),
    createdAt: formatDate(faker.date.past()),
  });
}

export function generateRandomUser(count: number): User[] {
  return Array.from({ length: count }, getRandomUser);
}

function getRandomClient(): Client {
  return clientSchema.parse({
    id: faker.number.int({ min: 1, max: 100 }),
    fullName: faker.person.fullName(),
    phoneNumber: faker.phone.number({ style: 'international' }),
    idNumber: faker.finance.accountNumber({ length: 8 }),
    dob: formatDate(faker.date.past()),
    gender: getRandomEnum(['MALE', 'FEMALE'] as const),
    active: getRandomBoolean(),
    branchName: faker.company.name(),
    assignedStaff: getRandomUser(),
    overpayment: faker.number.int({ min: 1, max: 10000 }),
    dueAmount: faker.number.int({ min: 1, max: 10000 }),
    createdBy: getRandomUser(),
    createdAt: formatDate(faker.date.past()),
  });
}

export function generateRandomClients(count: number): Client[] {
  return Array.from({ length: count }, getRandomClient);
}

function getRandomBranch(): Branch {
  return branchSchema.parse({
    id: faker.number.int({ min: 1, max: 100 }),
    branchName: faker.company.name(),
  });
}

export function generateRandomBranch(count: number): Branch[] {
  return Array.from({ length: count }, getRandomBranch);
}

function getRandomProduct(): Product {
  return productSchema.parse({
    id: faker.number.int({ min: 1, max: 100 }),
    branchName: faker.company.name(),
    loanAmount: faker.number.int({ min: 1, max: 100 }),
    repayAMount: faker.number.int({ min: 1, max: 100 }),
    interestAmount: faker.number.int({ min: 1, max: 100 }),
  });
}

export function generateRandomProduct(count: number): Product[] {
  return Array.from({ length: count }, getRandomProduct);
}

// Generate a random Loan
function generateRandomLoan(): Loan {
  const userSchema = getRandomUser();

  return loanSchema.parse({
    id: faker.number.int({ min: 1, max: 100 }),
    amount: faker.number.int({ min: 1, max: 20000 }),
    repayAmount: faker.number.int({ min: 1, max: 25000 }),
    client: getRandomClient(),
    loanOfficer: userSchema,
    loanPurpose: faker.lorem.lines(1),
    dueDate: formatDate(faker.date.future()),
    approvedBy: userSchema,
    disbursedOn: formatDate(faker.date.recent()),
    disbursedBy: userSchema,
    noOfInstallments: faker.number.int({ min: 1, max: 6 }),
    installmentsPeriod: faker.number.int({ min: 1, max: 30 }),
    status: getRandomEnum(['INACTIVE', 'ACTIVE', 'COMPLETED', 'DEFAULTED'] as const),
    processingFee: faker.number.int({ min: 1, max: 500 }),
    feePaid: faker.datatype.boolean(),
    paidAmount: faker.number.int({ min: 1, max: 25000 }),
    updatedBy: userSchema,
    createdBy: userSchema,
    createdAt: formatDate(faker.date.past()),
  });
}

// Generate a list of random loans
export function generateRandomLoans(count: number): Loan[] {
  return Array.from({ length: count }, generateRandomLoan);
}

function generateRandomInactiveLoan(): InactiveLoan {
  const userSchema = getRandomUser();

  return inactiveLoanSchema.parse({
    id: faker.number.int({ min: 1, max: 100 }),
    amount: faker.number.int({ min: 1, max: 20000 }),
    repayAmount: faker.number.int({ min: 1, max: 25000 }),
    client: getRandomClient(),
    loanOfficer: userSchema,
    approvedBy: userSchema,
    approvedOn: formatDate(faker.date.past()),
  });
}

// Generate a list of random inactive loans
export function generateRandomInactiveLoans(count: number): InactiveLoan[] {
  return Array.from({ length: count }, generateRandomInactiveLoan);
}
