<?php
$DB_HOST =  'database-mysql';
$DB_USER = 'case4';
$DB_PASSWORD = 'passdb_case4';
$DB_NAME = 'db_case4';

$conn = new mysqli($DB_HOST, $DB_USER, $DB_PASSWORD, $DB_NAME);

$jokes = [];
$error = null;

if ($conn->connect_error) {
    $error = "Connection failed: " . $conn->connect_error;
} else {
    $result = $conn->query("SELECT * FROM jokes ORDER BY created_at DESC");
    if ($result) {
        while ($row = $result->fetch_assoc()) {
            $jokes[] = $row;
        }
    }
    $conn->close();
}
?>
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chuck Norris Jokes</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            min-height: 100vh;
            padding: 20px;
            color: #fff;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            padding: 30px 0;
            margin-bottom: 30px;
        }

        header h1 {
            font-size: 2.5em;
            color: #ffd700;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
        }

        header p {
            color: #aaa;
            margin-top: 10px;
        }

        .jokes-container {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 20px;
        }

        .joke-card {
            background: rgba(255, 255, 255, 0.08);
            border-radius: 15px;
            padding: 25px;
            border-top: 4px solid #ffd700;
            transition: all 0.3s;
            display: flex;
            flex-direction: column;
            height: 100%;
        }

        .joke-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 30px rgba(255, 215, 0, 0.2);
            background: rgba(255, 255, 255, 0.12);
        }

        .joke-text {
            font-size: 1em;
            line-height: 1.6;
            color: #e0e0e0;
            flex-grow: 1;
        }

        .joke-footer {
            margin-top: 15px;
            font-size: 0.8em;
            color: #666;
            padding-top: 10px;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
        }

        .no-jokes {
            text-align: center;
            padding: 50px;
            color: #888;
            grid-column: 1 / -1;
        }

        .error {
            background: #ff4444;
            padding: 20px;
            border-radius: 10px;
            text-align: center;
        }

        footer {
            text-align: center;
            padding: 30px 0;
            margin-top: 40px;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            color: #888;
        }

        footer a {
            color: #ffd700;
            text-decoration: none;
            transition: all 0.3s;
        }

        footer a:hover {
            color: #ffed4a;
            text-decoration: underline;
        }

        @media (max-width: 768px) {
            .jokes-container {
                grid-template-columns: 1fr;
            }

            header h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>

<body>
    <div class="container">
        <header>
            <h1>Chuck Norris Jokes</h1>
            <p>Scraped from api.chucknorris.io every 10 seconds</p>
        </header>

        <?php if ($error): ?>
            <div class="error"><?= htmlspecialchars($error) ?></div>
        <?php else: ?>
            <div class="jokes-container">
                <?php if (empty($jokes)): ?>
                    <div class="no-jokes">
                        <h2>No jokes yet!</h2>
                        <p>Wait for the scrapper to fetch some jokes...</p>
                    </div>
                <?php else: ?>
                    <?php foreach ($jokes as $joke): ?>
                        <div class="joke-card">
                            <p class="joke-text"><?= htmlspecialchars($joke['joke_text']) ?></p>
                            <div class="joke-footer">Scraped at: <?= $joke['created_at'] ?></div>
                        </div>
                    <?php endforeach; ?>
                <?php endif; ?>
            </div>

            <footer>
                &copy; Cloud Computing <?= date('Y') ?> - <a href="https://github.com/LaptopMerah" target="_blank">LaptopMerah</a>
            </footer>
        <?php endif; ?>
    </div>

    <script>
        setTimeout(() => location.reload(), 15000);
    </script>
</body>

</html>
