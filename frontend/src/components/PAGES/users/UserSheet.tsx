import { useState } from 'react';
import { role } from '@/lib/types';
import { useAuth } from '@/hooks/useAuth';
import { useTable } from '@/hooks/useTable';
import { useQueryClient, useMutation, useQuery } from '@tanstack/react-query';
import updateUser from '@/services/updateUser';
import { updateUserType } from '@/lib/types';
import { toast } from 'react-toastify';
import getFormData from '@/services/getFormData';
import { format } from 'date-fns';
import VirtualizeddSelect from '../../UI/VisualizedSelect';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';

function UserSheet() {
	const [userRole, setUserRole] = useState<role | null>(null);
	const [branchId, setBranchId] = useState<number | null>(null);
	const { selectedRow, setSelectedRow } = useTable();
	const { decoded } = useAuth();

	const { data } = useQuery({
		queryKey: ['loans/form'],
		queryFn: () => getFormData(false, false, false, true),
		staleTime: 5 * 1000,
	});

	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: updateUser,
	});

	const onSave = () => {
		const values: updateUserType = {
			id: Number(selectedRow.id),
			role: userRole ? userRole : undefined,
			branchId: branchId ? branchId : undefined,
		};

		console.log(values);

		mutation.mutate(values, {
			onSuccess: () => {
				queryClient.invalidateQueries({ queryKey: ['users'] });
				toast.success('User Updated');
			},
			onError: (error) => {
				toast.error(error.message);
			},
			onSettled: () => mutation.reset(),
		});

		setUserRole(null);
		setBranchId(null);
		setSelectedRow(null);
	};

	const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
		id: (value) => (
			<Input
				readOnly
				placeholder={`LN${String(value).padStart(3, '0')}`}
				className="bg-gray-100 text-gray-500"
			/>
		),
		role: (value) => {
			return decoded?.role === role.ADMIN ? (
				<div>
					<Select
						defaultValue={value}
						// value={decoded?.role}
						onValueChange={(value) => {
							const tmp =
								value === 'ADMIN' ? role.ADMIN : role.AGENT;
							setUserRole(tmp);
						}}
					>
						<SelectTrigger className="w-[180px]">
							<SelectValue placeholder={value} />
						</SelectTrigger>
						<SelectContent>
							<SelectGroup>
								<SelectItem value={role.ADMIN}>
									ADMIN
								</SelectItem>
								<SelectItem value={role.AGENT}>
									AGENT
								</SelectItem>
							</SelectGroup>
						</SelectContent>
					</Select>
				</div>
			) : (
				<Input
					readOnly
					placeholder={value.toString()}
					className="bg-gray-100 text-gray-500"
				/>
			);
		},
		branchName: (value) => {
			return decoded?.role === role.ADMIN && data?.branch ? (
				<VirtualizeddSelect
					options={data.branch}
					placeholder={value}
					value={value}
					onChange={(id) => setBranchId(id)}
				/>
			) : (
				<Input
					readOnly
					placeholder={value.toString()}
					className="bg-gray-100 text-gray-500"
				/>
			);
		},
		createdAt: (value) => {
			return (
				<Input
					readOnly
					placeholder={format(value, 'PPP')}
					className="bg-gray-100 text-gray-500"
				/>
			);
		},
	};

	return (
		<Sheet
			open={!!selectedRow}
			onOpenChange={(open: boolean) => {
				if (!open) {
					setSelectedRow(null);
					setUserRole(null);
					setBranchId(null);
				}
			}}
		>
			<SheetContent className="overflow-auto custom-sheet-class">
				<SheetHeader>
					<SheetTitle>User Details</SheetTitle>
					<SheetDescription>Description goes here</SheetDescription>
				</SheetHeader>
				{selectedRow && (
					<div className="py-4">
						{Object.entries(selectedRow).map(([key, value]) => {
							if (key === 'createdBy' || key === 'updatedBy') {
								return;
							}
							if (fieldRenderers[key]) {
								return (
									<div
										key={key}
										className="grid grid-cols-[0.5fr_1fr] mb-4"
									>
										<span className="font-medium capitalize">
											{key}
										</span>
										{fieldRenderers[key](value)}
									</div>
								);
							}

							return (
								<div
									key={key}
									className="grid grid-cols-[0.5fr_1fr] mb-4"
								>
									<span className="font-medium capitalize">
										{key}
									</span>
									{typeof value === 'string' ||
									typeof value === 'number' ||
									typeof value === 'boolean' ? (
										<Input
											readOnly
											placeholder={value.toString()}
											className="bg-gray-100 text-gray-500"
										/>
									) : (
										JSON.stringify(value)
									)}
								</div>
							);
						})}
						<Button size="lg" onClick={onSave} className="mt-8">
							Save
						</Button>
					</div>
				)}
			</SheetContent>
		</Sheet>
	);
}

export default UserSheet;
