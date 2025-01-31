import { Product } from '../PAGES/products/schema';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

interface ProductCardDisplayProps {
	product: Product;
}

function ProductCardDisplay({ product }: ProductCardDisplayProps) {
	return (
		<Card className="p-2" data-user-id={product.id}>
			<CardHeader className="p-0 mb-2">
				<CardTitle className="text-lg font-semibold">
					{`${product.loanAmount} - ${product.branchName}`}
				</CardTitle>
			</CardHeader>
			<div className="flex gap-x-2 justify-between">
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Repay Amount
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{product.repayAmount}</p>
					</CardContent>
				</div>
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Interest
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{product.interestAmount}</p>
					</CardContent>
				</div>
			</div>
		</Card>
	);
}

export default ProductCardDisplay;
