import {generateRandomBranch, generateRandomClients, generateRandomInactiveLoans, generateRandomProduct, generateRandomLoans, generateRandomUser } from "@/lib/generator";

const branches = generateRandomBranch(30);
const customers = generateRandomClients(30);
const generatedInactiveLoans = generateRandomInactiveLoans(13);
const generatedLoans = generateRandomLoans(30);
const products = generateRandomProduct(30);
const users = generateRandomUser(30);

function GetDataTest() {
    const data = {
        branches,
        customers,
        inactiveLoans: generatedInactiveLoans,
        loans: generatedLoans,
        products,
        users,
      };
    
      return (
        <div>
          <pre>{JSON.stringify(data, null, 2)}</pre>
        </div>
      );
}

export default GetDataTest