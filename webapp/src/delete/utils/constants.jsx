// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable no-magic-numbers */

import audioIcon from '../images/icons/audio.svg';
import codeIcon from '../images/icons/code.svg';
import excelIcon from '../images/icons/excel.svg';
import genericIcon from '../images/icons/generic.svg';
import patchIcon from '../images/icons/patch.png';
import pdfIcon from '../images/icons/pdf.svg';
import pptIcon from '../images/icons/ppt.svg';
import videoIcon from '../images/icons/video.svg';
import wordIcon from '../images/icons/word.svg';

import {t} from './i18n';

export const FileTypes = {
  IMAGE: 'image',
  AUDIO: 'audio',
  VIDEO: 'video',
  SPREADSHEET: 'spreadsheet',
  CODE: 'code',
  WORD: 'word',
  PRESENTATION: 'presentation',
  PDF: 'pdf',
  PATCH: 'patch',
  SVG: 'svg',
  OTHER: 'other',
};


export const Constants = {
  IMAGE_TYPE_GIF: 'gif',
  IMAGE_TYPES: ['jpg', 'gif', 'bmp', 'png', 'jpeg', 'tiff', 'tif'],
  AUDIO_TYPES: ['mp3', 'wav', 'wma', 'm4a', 'flac', 'aac', 'ogg', 'm4r'],
  VIDEO_TYPES: ['mp4', 'avi', 'webm', 'mkv', 'wmv', 'mpg', 'mov', 'flv'],
  PRESENTATION_TYPES: ['ppt', 'pptx'],
  SPREADSHEET_TYPES: ['xlsx', 'csv'],
  WORD_TYPES: ['doc', 'docx'],
  CODE_TYPES: ['applescript', 'as', 'atom', 'bas', 'bash', 'boot', 'c', 'c++', 'cake', 'cc', 'cjsx', 'cl2', 'clj', 'cljc', 'cljs', 'cljs.hl', 'cljscm', 'cljx', '_coffee', 'coffee', 'cpp', 'cs', 'csharp', 'cson', 'css', 'd', 'dart', 'delphi', 'dfm', 'di', 'diff', 'django', 'docker', 'dockerfile', 'dpr', 'erl', 'ex', 'exs', 'f90', 'f95', 'freepascal', 'fs', 'fsharp', 'gcode', 'gemspec', 'go', 'groovy', 'gyp', 'h', 'h++', 'handlebars', 'hbs', 'hic', 'hpp', 'hs', 'html', 'html.handlebars', 'html.hbs', 'hx', 'iced', 'irb', 'java', 'jinja', 'jl', 'js', 'json', 'jsp', 'jsx', 'kt', 'ktm', 'kts', 'lazarus', 'less', 'lfm', 'lisp', 'log', 'lpr', 'lua', 'm', 'mak', 'matlab', 'md', 'mk', 'mkd', 'mkdown', 'ml', 'mm', 'nc', 'obj-c', 'objc', 'osascript', 'pas', 'pascal', 'perl', 'php', 'php3', 'php4', 'php5', 'php6', 'pl', 'plist', 'podspec', 'pp', 'ps', 'ps1', 'py', 'r', 'rb', 'rs', 'rss', 'ruby', 'scala', 'scm', 'scpt', 'scss', 'sh', 'sld', 'sql', 'st', 'styl', 'swift', 'tex', 'thor', 'txt', 'v', 'vb', 'vbnet', 'vbs', 'veo', 'xhtml', 'xml', 'xsl', 'yaml', 'zsh'],
  PDF_TYPES: ['pdf'],
  PATCH_TYPES: ['patch'],
  SVG_TYPES: ['svg'],
  ICON_FROM_TYPE: {
      audio: audioIcon,
      video: videoIcon,
      spreadsheet: excelIcon,
      presentation: pptIcon,
      pdf: pdfIcon,
      code: codeIcon,
      word: wordIcon,
      patch: patchIcon,
      other: genericIcon,
  },
};

t('suggestion.mention.channels');
t('suggestion.mention.morechannels');
t('suggestion.mention.unread.channels');
t('suggestion.mention.members');
t('suggestion.mention.moremembers');
t('suggestion.mention.nonmembers');
t('suggestion.mention.special');
t('suggestion.archive');

export default Constants;
