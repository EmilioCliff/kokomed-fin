import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { useContext, useState } from 'react';
import { TableContext } from '@/context/TableContext';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
  SelectSeparator,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Calendar } from '@/components/ui/calendar';
import { isUser, isClient } from '@/lib/validators';
import UserCardDisplay from '@/components/UI/UserCardDisplay';
import ClientCardDisplay from '@/components/UI/ClientCardDisplay';
import { CalendarIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';
import { updateLoan } from '@/services/updateLoan';
import { updateLoanType, loanStatus } from '@/lib/types';

function LoanSheet() {
  const [status, setStatus] = useState<loanStatus | null>(null);
  const [disbursedDate, setDisbursedDate] = useState<string | null>(null);
  const [feePaid, setFeePaid] = useState<boolean | null>(null);
  const { selectedRow, setSelectedRow } = useContext(TableContext);

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: updateLoan,
  });

  const onSave = () => {
    const values: updateLoanType = {
      id: Number(selectedRow.id),
      status: status ? status : undefined,
      disburseDate: disbursedDate ? disbursedDate : undefined,
      feePaid: feePaid ? feePaid : undefined,
    };

    console.log(values);

    mutation.mutate(values, {
      onSuccess: (data) => {
        console.log(data);
        // we will invalidate but generally we need to add the updated loan into the state
        queryClient.invalidateQueries({ queryKey: 'loans' });
      },
      onError: (error) => {
        console.log(error);
        // show to toastify
      },
    });

    setStatus(null);
    setDisbursedDate(null);
    setFeePaid(null);
    queryClient.invalidateQueries({
      queryKey: ['loans'],
    });
  };

  const isAdmin = true;

  const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
    id: (value) => (
      <Input
        readOnly
        placeholder={`LN${String(value).padStart(3, '0')}`}
        className="bg-gray-100 text-gray-500"
      />
    ),
    client: (value) => <ClientCardDisplay client={isClient(value)} />,
    loanOfficer: (value) => <UserCardDisplay user={isUser(value)} />,
    loanPurpose: (value) => <>{value}</>,
    approvedBy: (value) => <UserCardDisplay user={isUser(value)} />,
    disbursedBy: (value) => <UserCardDisplay user={isUser(value)} />,
    status: (value: string) => {
      return value === 'INACTIVE' && isAdmin ? (
        <div>
          <Select
            defaultValue="INACTIVE"
            onValueChange={(value) => {
              const tmp = value === 'INACTIVE' ? loanStatus.INACTIVE : loanStatus.ACTIVE;
              setStatus(tmp);
            }}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="INACTIVE" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="INACTIVE">INACTIVE</SelectItem>
                <SelectItem value="ACTIVE">ACTIVE</SelectItem>
                <SelectSeparator />
                <SelectLabel>Disbursed Date</SelectLabel>
                <Popover>
                  <PopoverTrigger asChild>
                    <Button
                      variant={'outline'}
                      className={cn(
                        'w-[240px] pl-3 text-left font-normal',
                        !disbursedDate && 'text-muted-foreground'
                      )}
                    >
                      {disbursedDate ? (
                        format(disbursedDate, 'PPP')
                      ) : (
                        <span>Pick a date</span>
                      )}
                      <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={disbursedDate ? new Date(disbursedDate) : undefined}
                      onSelect={(date) => setDisbursedDate(format(date!, 'yyyy-MM-dd'))}
                      disabled={(date) =>
                        date > new Date() || date < new Date('1900-01-01')
                      }
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
              </SelectGroup>
            </SelectContent>
          </Select>
          {disbursedDate && (
            <p>Disbursed Date: {format(new Date(disbursedDate), 'PPP')}</p>
          )}
        </div>
      ) : (
        <Input
          readOnly
          placeholder={value.toString()}
          className="bg-gray-100 text-gray-500"
        />
      );
    },
    feePaid: (value) => {
      return value === false ? (
        <Select defaultValue="NO" onValueChange={(value) => setFeePaid(value === 'YES')}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="NO" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem value="YES">YES</SelectItem>
              <SelectItem value="NO">NO</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      ) : (
        <Input readOnly placeholder="YES" className="bg-gray-100 text-gray-500" />
      );
    },
    updatedBy: () => <></>,
    createdBy: () => <></>,
  };

  return (
    <Sheet
      open={!!selectedRow}
      onOpenChange={(open: boolean) => {
        if (!open) {
          setSelectedRow(null);
          setDisbursedDate(null);
          setFeePaid(null);
        }
      }}
    >
      <SheetContent className="overflow-auto custom-sheet-class">
        <SheetHeader>
          <SheetTitle>Loan Details</SheetTitle>
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
                  <div key={key} className="grid grid-cols-[0.5fr_1fr] mb-4">
                    <span className="font-medium capitalize">{key}</span>
                    {fieldRenderers[key](value)}
                  </div>
                );
              }

              return (
                <div key={key} className="grid grid-cols-[0.5fr_1fr] mb-4">
                  <span className="font-medium capitalize">{key}</span>
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

export default LoanSheet;
