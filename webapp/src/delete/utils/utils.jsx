import Constants, {FileTypes} from './constants.jsx';

// Converts a file size in bytes into a human-readable string of the form '123MB'.
export function fileSizeToString(bytes) {
  // it's unlikely that we'll have files bigger than this
  if (bytes > 1024 * 1024 * 1024 * 1024) {
      return Math.floor(bytes / (1024 * 1024 * 1024 * 1024)) + 'TB';
  } else if (bytes > 1024 * 1024 * 1024) {
      return Math.floor(bytes / (1024 * 1024 * 1024)) + 'GB';
  } else if (bytes > 1024 * 1024) {
      return Math.floor(bytes / (1024 * 1024)) + 'MB';
  } else if (bytes > 1024) {
      return Math.floor(bytes / 1024) + 'KB';
  }

  return bytes + 'B';
}

export function getFileIconPath(fileInfo) {
  const fileType = getFileType(fileInfo.extension);

  var icon;
  if (fileType in Constants.ICON_FROM_TYPE) {
      icon = Constants.ICON_FROM_TYPE[fileType];
  } else {
      icon = Constants.ICON_FROM_TYPE.other;
  }

  return icon;
}

export const getFileType = (extin) => {
  const ext = removeQuerystringOrHash(extin.toLowerCase());

  if (Constants.IMAGE_TYPES.indexOf(ext) > -1) {
      return FileTypes.IMAGE;
  }

  if (Constants.AUDIO_TYPES.indexOf(ext) > -1) {
      return FileTypes.AUDIO;
  }

  if (Constants.VIDEO_TYPES.indexOf(ext) > -1) {
      return FileTypes.VIDEO;
  }

  if (Constants.SPREADSHEET_TYPES.indexOf(ext) > -1) {
      return FileTypes.SPREADSHEET;
  }

  if (Constants.CODE_TYPES.indexOf(ext) > -1) {
      return FileTypes.CODE;
  }

  if (Constants.WORD_TYPES.indexOf(ext) > -1) {
      return FileTypes.WORD;
  }

  if (Constants.PRESENTATION_TYPES.indexOf(ext) > -1) {
      return FileTypes.PRESENTATION;
  }

  if (Constants.PDF_TYPES.indexOf(ext) > -1) {
      return FileTypes.PDF;
  }

  if (Constants.PATCH_TYPES.indexOf(ext) > -1) {
      return FileTypes.PATCH;
  }

  if (Constants.SVG_TYPES.indexOf(ext) > -1) {
      return FileTypes.SVG;
  }

  return FileTypes.OTHER;
};

const removeQuerystringOrHash = (extin) => {
  return extin.split(/[?#]/)[0];
};
