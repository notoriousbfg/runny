import * as vscode from 'vscode'

export function activate(context: vscode.ExtensionContext) {
    vscode.languages.registerDocumentFormattingEditProvider('rny', {
        provideDocumentFormattingEdits(document: vscode.TextDocument): vscode.TextEdit[] {
            const edits: vscode.TextEdit[] = []
            for (let i = 0; i < document.lineCount; i++) {
                const line = document.lineAt(i)
                const trimmed = line.text.trim()

                if (trimmed.startsWith("}") || trimmed.startsWith("run") || trimmed.startsWith("var")) {
                    const indentLevel = line.firstNonWhitespaceCharacterIndex
                    const newText = line.text.trimStart()
                    edits.push(vscode.TextEdit.replace(line.range, ' '.repeat(indentLevel) + newText))
                }
            }
            return edits
        }
    })
}
