import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PieChart, Pie, Cell, ResponsiveContainer, Legend } from "recharts";

const data = [
	{ name: "Active", value: 20 },
	{ name: "Inactive", value: 15 },
	{ name: "Completed", value: 40 },
	{ name: "Disbursed", value: 25 },
];

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"];

function LoanStatusChart() {
	return (
		<Card className='col-span-1'>
			<CardHeader>
				<CardTitle>Loan Status</CardTitle>
			</CardHeader>
			<CardContent>
				<ResponsiveContainer width='100%' height={300}>
					<PieChart>
						<Pie
							data={data}
							cx='50%'
							cy='50%'
							labelLine={false}
							outerRadius={80}
							fill='#8884d8'
							dataKey='value'
						>
							{data.map((entry, index) => (
								<Cell
									key={`cell-${index}`}
									fill={COLORS[index % COLORS.length]}
								/>
							))}
						</Pie>
						<Legend />
					</PieChart>
				</ResponsiveContainer>
			</CardContent>
		</Card>
	);
}

export default LoanStatusChart;
