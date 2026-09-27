using System;
using System.Diagnostics;
using System.IO;
using System.Threading;
using System.Windows.Forms;
using DiscordRPC;
using DiscordRPC.Logging;
using Button = DiscordRPC.Button;

namespace RealmOfChaos
{
    internal static class Program
    {
        private const string CLIENT_ID = "YOUR_DISCORD_APPLICATION_ID";

        [STAThread]
        static void Main(string[] args)
        {
            string gamePath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "element", "elementclient.exe");

            // Fallback checking
            if (!File.Exists(gamePath))
            {
                gamePath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "elementclient.exe");
            }

            if (!File.Exists(gamePath))
            {
                MessageBox.Show(
                    "File 'element/elementclient.exe' tidak ditemukan!\n\nPastikan launcher ini diletakkan di folder utama game Perfect World Anda (sejajar dengan folder 'element').",
                    "Realm of Chaos - Error",
                    MessageBoxButtons.OK,
                    MessageBoxIcon.Warning
                );
                return;
            }

            // Initialize Discord RPC Client
            DiscordRpcClient client = new DiscordRpcClient(CLIENT_ID)
            {
                Logger = new ConsoleLogger() { Level = LogLevel.Warning }
            };

            client.Initialize();

            // Set Rich Presence Activity
            client.SetPresence(new RichPresence()
            {
                Details = "Perfect World v1.4.6",
                State = "Playing on Realm of Chaos",
                Timestamps = Timestamps.Now,
                Assets = new Assets()
                {
                    LargeImageKey = "logo_roc",
                    LargeImageText = "Realm of Chaos - Sirens of War",
                    SmallImageKey = "pwi",
                    SmallImageText = "v1.4.6 build 2305"
                },
                Buttons = new Button[]
                {
                    new Button() { Label = "🌐 Website", Url = "https://your-server-website.com" },
                    new Button() { Label = "💬 Discord Server", Url = "https://discord.gg/your-discord" }
                }
            });

            // Start game client
            ProcessStartInfo psi = new ProcessStartInfo()
            {
                FileName = gamePath,
                Arguments = "game:cpw console:1",
                WorkingDirectory = Path.GetDirectoryName(gamePath) ?? AppDomain.CurrentDomain.BaseDirectory,
                UseShellExecute = true
            };

            try
            {
                using (Process? gameProcess = Process.Start(psi))
                {
                    if (gameProcess != null)
                    {
                        // Wait until the player exits the game
                        gameProcess.WaitForExit();
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show("Gagal memulai game:\n" + ex.Message, "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
            finally
            {
                // Clear presence and dispose
                client.ClearPresence();
                client.Dispose();
            }
        }
    }
}
