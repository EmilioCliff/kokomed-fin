import api from '@/API/api';
import { commonresponse } from '@/lib/types';
import { ReportFilter } from '@/lib/types';

const generateReport = async (data: ReportFilter) => {
	try {
		const response = await api.post(`/report?format=${data.format}`, data, {
			responseType: 'blob',
		});

		const contentDisposition = response.headers['content-disposition'];
		let fileName = `report_${Date.now()}.${
			data.format === 'excel' ? 'xlsx' : 'pdf'
		}`;

		if (contentDisposition) {
			const fileNameMatch = contentDisposition.match(
				/filename="?([^";]+)"?/,
			);
			if (fileNameMatch && fileNameMatch[1]) {
				fileName = fileNameMatch[1];
			}
		}

		const blob = new Blob([response.data]);
		const fileURL = window.URL.createObjectURL(blob);
		const fileLink = document.createElement('a');

		fileLink.href = fileURL;
		fileLink.setAttribute('download', fileName);
		document.body.appendChild(fileLink);
		fileLink.click();

		document.body.removeChild(fileLink);
		URL.revokeObjectURL(fileURL);

		return response;
	} catch (error: any) {
		if (error.response) {
			throw new Error(error.response.data.message);
		}

		throw new Error(error.message);
	}
};

export default generateReport;
