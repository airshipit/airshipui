import { NgxMonacoEditorConfig } from 'ngx-monaco-editor';

const monacoConfig: NgxMonacoEditorConfig = {
    onMonacoLoad: () => {
        monaco.editor.defineTheme('airshipTheme', {
        base: 'vs',
        inherit: true,
        rules: [],
            colors: {
                'editor.background': '#f5f5f5'
            }
        });
    }
};

export default monacoConfig;
