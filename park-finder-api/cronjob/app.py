from src import create_app


def main():
    print("Start scheduler...")
    app = create_app()
    app.run(host="0.0.0.0", port=4200, debug=True)


if __name__ == "__main__":
    main()
    