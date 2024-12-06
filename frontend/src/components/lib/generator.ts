import { faker } from "@faker-js/faker";
import { z } from "zod";
import {
	loanSchema,
	Loan,
	User,
	userSchema,
	InactiveLoan,
	Client,
	inactiveLoanSchema,
} from "../../data/schema"; // Adjust path to your schema file
import InactiveLoans from "../InactiveLoans";

function getRandomEnum<T extends string>(enumValues: T[]): T {
	return enumValues[Math.floor(Math.random() * enumValues.length)];
}

function getRandomBoolean(): boolean {
	return Math.random() < 0.5;
}

function getRandomUser(): User {
	return {
		id: faker.number.int({ min: 1, max: 100 }),
		fullName: faker.person.fullName(),
		phoneNumber: faker.phone.number({ style: "international" }),
		email: faker.internet.email(),
		role: getRandomEnum(["ADMIN", "AGENT"] as const),
	};
}

function getRandomClient(): Client {
	return {
		id: faker.number.int({ min: 1, max: 100 }),
		fullName: faker.person.fullName(),
		phoneNumber: faker.phone.number({ style: "international" }),
		active: faker.datatype.boolean(),
		assignedStaff: getRandomUser(),
		overpayment: faker.number.int({ min: 1, max: 10000 }),
	};
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
		dueDate: faker.date.future().toISOString(),
		approvedBy: userSchema,
		disbursedOn: faker.date.recent().toISOString(),
		disbursedBy: userSchema,
		noOfInstallments: faker.number.int({ min: 1, max: 6 }),
		installmentsPeriod: faker.number.int({ min: 1, max: 30 }),
		status: getRandomEnum([
			"INACTIVE",
			"ACTIVE",
			"COMPLETED",
			"DEFAULTED",
		] as const),
		processingFee: faker.number.int({ min: 1, max: 500 }),
		feePaid: faker.datatype.boolean(),
		paidAmount: faker.number.int({ min: 1, max: 25000 }),
		updatedBy: userSchema,
		createdBy: userSchema,
		createdAt: faker.date.past().toISOString(),
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
		approvedOn: faker.date.past().toISOString(),
	});
}

// Generate a list of random inactive loans
export function generateRandomInactiveLoans(count: number): InactiveLoan[] {
	return Array.from({ length: count }, generateRandomInactiveLoan);
}
