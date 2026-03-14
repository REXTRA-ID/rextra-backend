param(
    [string]$msg = 'Tugas Selesai!',
    [string]$title = 'Gemini CLI Notification',
    [string]$type = 'Info' # Info, Error, Choice
)

$soundPath = 'C:\Users\USER\Downloads\jokowi-kaget.mp3'
Add-Type -AssemblyName PresentationCore
$player = New-Object system.windows.media.mediaplayer
$player.open($soundPath)
$player.Play()

Add-Type -AssemblyName System.Windows.Forms
$icon = [System.Windows.Forms.MessageBoxIcon]::Information
$buttons = [System.Windows.Forms.MessageBoxButtons]::OK

if ($type -eq 'Error') {
    $icon = [System.Windows.Forms.MessageBoxIcon]::Error
} elseif ($type -eq 'Choice') {
    $icon = [System.Windows.Forms.MessageBoxIcon]::Question
    $buttons = [System.Windows.Forms.MessageBoxButtons]::YesNo
}

$result = [System.Windows.Forms.MessageBox]::Show($msg, $title, $buttons, $icon)
return $result
