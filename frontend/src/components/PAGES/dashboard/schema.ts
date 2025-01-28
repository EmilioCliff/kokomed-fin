import { z } from 'zod';
import { clientSchema } from '../customers/schema';
import { userSchema } from '../users/schema';

export const inactiveLoanSchema = z.object({
  id: z.number(),
  amount: z.number(),
  repayAmount: z.number(),
  client: clientSchema,
  loanOfficer: userSchema,
  approvedBy: userSchema,
  approvedOn: z
    .string()
    .optional()
    .refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
      message: 'Invalid date string!',
    }),
});

export type InactiveLoan = z.infer<typeof inactiveLoanSchema>;

// recent payments
export const recentPaymentSchema = z.object({
  id: z.number(),
  borrower: z.string(),
  amount: z.number(),
  date: z.string(),
});

export type RecentPayment = z.infer<typeof recentPaymentSchema>;

// dashboard widgets
export const widgetSchema = z.object({
  title: z.string(),
  mainAmount: z.number().optional(),
  active: z.number().optional(),
  activeTitle: z.string(),
  closed: z.number().optional(),
  closedTitle: z.string(),
  currency: z.string().optional(),
});

export type Widget = z.infer<typeof widgetSchema>;

export const dashboardSchema = z.object({
  widgets: z.array(widgetSchema),
  recentPayments: z.array(recentPaymentSchema),
  inactiveLoans: z.array(inactiveLoanSchema),
});

export type DashboardData = z.infer<typeof dashboardSchema>;
