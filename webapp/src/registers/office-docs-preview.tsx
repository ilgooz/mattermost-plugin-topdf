import { Store } from 'redux';
import { Provider } from 'react-redux';
import { FileInfo } from 'mattermost-redux/types/files'

import { id } from '../manifest'
import PDFPreview from '../delete/components/pdf_preview';

export default (store: Store) => {
  return ({ fileInfo }: { info: FileInfo }) => {
    const previewURL = `/plugins/${id}/files/${fileInfo.id}`;

    return (
      <Provider store={store}>
        <PDFPreview fileInfo={fileInfo} fileUrl={previewURL}></PDFPreview>
      </Provider>
    )
  }
}

const supportedExtensions = [ "doc", "docx", "odt", "xls", "xlsx", "ods", "ppt", "pptx", "odp" ];

export const isOfficeDocument = (extension: string) => supportedExtensions.includes(extension);
